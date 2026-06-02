#!/usr/bin/env python3
"""
周报生成辅助脚本 - Weekly Report Generator Helper

提供一些辅助功能来帮助生成周报：
1. 日期范围自动计算（本周、上周等）
2. 数据格式校验和转化
3. 周报模板加载
"""

import json
from datetime import datetime, timedelta


def get_week_range(date_str=None):
    """获取某一天所在周的周一和周日日期"""
    if date_str:
        date = datetime.strptime(date_str, "%Y-%m-%d")
    else:
        date = datetime.now()
    
    monday = date - timedelta(days=date.weekday())
    sunday = monday + timedelta(days=6)
    
    return {
        "monday": monday.strftime("%Y-%m-%d"),
        "sunday": sunday.strftime("%Y-%m-%d"),
        "week_range": f"{monday.strftime('%m/%d')} - {sunday.strftime('%m/%d')}",
        "year_week": f"{monday.isocalendar()[0]}W{monday.isocalendar()[1]}"
    }


def format_percentage(current, target):
    """计算完成率并格式化"""
    if target == 0:
        return "N/A"
    rate = (current / target) * 100
    return f"{rate:.1f}%"


def format_change(current, previous):
    """计算环比变化并格式化"""
    if previous == 0:
        return "N/A"
    change = ((current - previous) / previous) * 100
    symbol = "↑" if change > 0 else ("↓" if change < 0 else "→")
    return f"{symbol} {abs(change):.1f}%"


def generate_status_badge(progress):
    """根据进度生成状态标记"""
    if progress >= 100:
        return "✅ 已完成"
    elif progress >= 80:
        return "🟢 收尾中"
    elif progress >= 50:
        return "🔄 进行中"
    elif progress >= 20:
        return "🟡 启动中"
    else:
        return "⏳ 待启动"


def highlight_keywords(text):
    """识别并标记文本中的关键信息"""
    keywords = []
    
    # 数字相关
    import re
    numbers = re.findall(r'\d+[%倍%]?', text)
    if numbers:
        keywords.append({"type": "data", "values": numbers})
    
    # 程度词
    intensity_words = ["重大", "关键", "核心", "重要", "紧急", "突破", "首次"]
    for word in intensity_words:
        if word in text:
            keywords.append({"type": "intensity", "word": word})
            break
    
    return keywords


def suggest_template(role_type):
    """根据岗位类型推荐模板"""
    templates = {
        "研发": "研发工程师模板",
        "开发": "研发工程师模板",
        "后端": "研发工程师模板",
        "前端": "研发工程师模板",
        "算法": "研发工程师模板",
        "测试": "研发工程师模板",
        "产品": "产品经理模板",
        "产品经理": "产品经理模板",
        "运营": "运营类模板",
        "市场": "运营类模板",
        "销售": "销售/商务类模板",
        "商务": "销售/商务类模板",
        "管理": "管理岗模板",
        "经理": "管理岗模板",
        "总监": "管理岗模板",
        "设计": "通用型模板（可用研发模板调整）",
        "UI": "通用型模板",
        "UX": "通用型模板",
    }
    
    for key, template in templates.items():
        if key in role_type:
            return template
    return "通用型模板"


if __name__ == "__main__":
    # 示例用法
    week = get_week_range()
    print(f"本周范围: {week['week_range']}")
    print(f"本周编号: {week['year_week']}")
    print(f"建议模板: {suggest_template('后端开发')}")
