package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ESHit 表示一次搜索命中的单条文档。
type ESHit struct {
	ID     string          `json:"id"`
	Score  float64         `json:"score"`
	Source json.RawMessage `json:"source"`
}

// ESSearchResult 表示一次搜索的聚合结果。
type ESSearchResult struct {
	Total int     `json:"total"`
	Hits  []ESHit `json:"hits"`
}

// EnsureESIndex 幂等地创建索引；当 mapping 为空时仅在索引不存在时创建空索引。
// mapping 形如 {"mappings":{"properties":{...}}}。
func (s *Store) EnsureESIndex(ctx context.Context, index string, mapping map[string]any) error {
	exists, err := s.ESIndexExists(ctx, index)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	req := esapi.IndicesCreateRequest{Index: index}
	if len(mapping) > 0 {
		body, err := json.Marshal(mapping)
		if err != nil {
			return fmt.Errorf("marshal index mapping: %w", err)
		}
		req.Body = bytes.NewReader(body)
	}
	res, err := req.Do(ctx, s.ES)
	if err != nil {
		return fmt.Errorf("create index %s: %w", index, err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("create index %s: %s", index, res.String())
	}
	return nil
}

// ESIndexExists 判断索引是否存在。
func (s *Store) ESIndexExists(ctx context.Context, index string) (bool, error) {
	res, err := s.ES.Indices.Exists([]string{index}, s.ES.Indices.Exists.WithContext(ctx))
	if err != nil {
		return false, fmt.Errorf("check index %s: %w", index, err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		return true, nil
	}
	if res.StatusCode == 404 {
		return false, nil
	}
	return false, fmt.Errorf("check index %s: %s", index, res.String())
}

// IndexDocument 插入或整体覆盖一条文档（指定 ID 时为 upsert）。
func (s *Store) IndexDocument(ctx context.Context, index, id string, doc any) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "true", // 写入后立即可查；高吞吐场景可改为 "false"
	}
	res, err := req.Do(ctx, s.ES)
	if err != nil {
		return fmt.Errorf("index document: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("index document: %s", res.String())
	}
	return nil
}

// GetDocument 按 ID 获取文档，并反序列化到 dest。不存在时返回 found=false。
func (s *Store) GetDocument(ctx context.Context, index, id string, dest any) (bool, error) {
	res, err := s.ES.Get(index, id, s.ES.Get.WithContext(ctx))
	if err != nil {
		return false, fmt.Errorf("get document: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return false, nil
	}
	if res.IsError() {
		return false, fmt.Errorf("get document: %s", res.String())
	}

	var envelope struct {
		Source json.RawMessage `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return false, fmt.Errorf("decode get response: %w", err)
	}
	if dest != nil {
		if err := json.Unmarshal(envelope.Source, dest); err != nil {
			return false, fmt.Errorf("unmarshal source: %w", err)
		}
	}
	return true, nil
}

// SearchDocuments 执行一次 Query DSL 搜索。query 形如 {"query":{...},"size":10}。
func (s *Store) SearchDocuments(ctx context.Context, index string, query map[string]any) (ESSearchResult, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return ESSearchResult{}, fmt.Errorf("marshal query: %w", err)
	}
	res, err := s.ES.Search(
		s.ES.Search.WithContext(ctx),
		s.ES.Search.WithIndex(index),
		s.ES.Search.WithBody(bytes.NewReader(body)),
		s.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return ESSearchResult{}, fmt.Errorf("search documents: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		raw, _ := io.ReadAll(res.Body)
		return ESSearchResult{}, fmt.Errorf("search documents: %s", string(raw))
	}

	var envelope struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string          `json:"_id"`
				Score  float64         `json:"_score"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return ESSearchResult{}, fmt.Errorf("decode search response: %w", err)
	}

	result := ESSearchResult{Total: envelope.Hits.Total.Value}
	for _, h := range envelope.Hits.Hits {
		result.Hits = append(result.Hits, ESHit{ID: h.ID, Score: h.Score, Source: h.Source})
	}
	return result, nil
}

// MatchDocuments 是按单字段精确/全文匹配的便捷查询封装。
func (s *Store) MatchDocuments(ctx context.Context, index, field string, value any, size int) (ESSearchResult, error) {
	if size <= 0 {
		size = 10
	}
	query := map[string]any{
		"size": size,
		"query": map[string]any{
			"match": map[string]any{field: value},
		},
	}
	return s.SearchDocuments(ctx, index, query)
}

// DeleteDocument 按 ID 删除文档；文档不存在视为成功。
func (s *Store) DeleteDocument(ctx context.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}
	res, err := req.Do(ctx, s.ES)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return nil
	}
	if res.IsError() {
		return fmt.Errorf("delete document: %s", res.String())
	}
	return nil
}

// BulkIndexDocuments 批量插入/覆盖文档。docs 的 key 为文档 ID，value 为文档内容。
func (s *Store) BulkIndexDocuments(ctx context.Context, index string, docs map[string]any) error {
	if len(docs) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for id, doc := range docs {
		meta := map[string]any{"index": map[string]any{"_index": index, "_id": id}}
		metaLine, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("marshal bulk meta: %w", err)
		}
		docLine, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("marshal bulk doc: %w", err)
		}
		buf.Write(metaLine)
		buf.WriteByte('\n')
		buf.Write(docLine)
		buf.WriteByte('\n')
	}

	res, err := s.ES.Bulk(
		strings.NewReader(buf.String()),
		s.ES.Bulk.WithContext(ctx),
		s.ES.Bulk.WithIndex(index),
		s.ES.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("bulk index: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("bulk index: %s", res.String())
	}

	var envelope struct {
		Errors bool `json:"errors"`
		Items  []map[string]struct {
			Status int             `json:"status"`
			Error  json.RawMessage `json:"error"`
		} `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode bulk response: %w", err)
	}
	if envelope.Errors {
		for _, item := range envelope.Items {
			for _, op := range item {
				if op.Error != nil {
					return fmt.Errorf("bulk index partial failure: %s", string(op.Error))
				}
			}
		}
	}
	return nil
}
