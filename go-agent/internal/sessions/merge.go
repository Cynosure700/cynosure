package sessions

func MergeSkillLoaders(loaders ...*SkillLoader) *SkillLoader {
	merged := NewSkillLoader()
	entries := make(map[string]*SkillEntry)
	for _, loader := range loaders {
		if loader == nil {
			continue
		}
		for name, entry := range loader.Entries() {
			entries[name] = entry
		}
	}
	merged.LoadFromEntries(entries)
	return merged
}

func cloneSkillEntry(entry *SkillEntry) *SkillEntry {
	if entry == nil {
		return nil
	}
	meta := make(map[string]string, len(entry.Meta))
	for key, value := range entry.Meta {
		meta[key] = value
	}
	return &SkillEntry{
		Meta: meta,
		Body: entry.Body,
		Path: entry.Path,
	}
}
