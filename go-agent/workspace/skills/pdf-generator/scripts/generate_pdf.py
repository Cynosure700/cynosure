#!/usr/bin/env python3
"""
PDF Generator Utility Script
A reusable helper that handles common PDF generation patterns.

Usage:
    python generate_pdf.py fpdf2 --output out.pdf --title "My Doc" --content text.txt
    python generate_pdf.py reportlab --output report.pdf --title "Report" --data data.json
    python generate_pdf.py weasyprint --output out.pdf --html content.html
    python generate_pdf.py markdown --output doc.pdf --markdown doc.md
"""

import argparse
import json
import os
import sys
from pathlib import Path


def cmd_fpdf2(args):
    """Generate a PDF using fpdf2 from text or structured data."""
    from fpdf import FPDF

    pdf = FPDF()
    pdf.set_auto_page_break(auto=True, margin=15)
    pdf.add_page()

    if args.title:
        pdf.set_font("Helvetica", "B", 20)
        pdf.cell(0, 15, args.title, align="C", new_x="LMARGIN", new_y="NEXT_LINE")
        pdf.ln(5)

    if args.content:
        pdf.set_font("Helvetica", size=12)
        content_path = Path(args.content)
        if content_path.exists():
            text = content_path.read_text(encoding="utf-8")
        else:
            text = args.content
        for paragraph in text.split("\n\n"):
            pdf.multi_cell(0, 8, paragraph.strip())
            pdf.ln(2)

    if args.data:
        data_path = Path(args.data)
        if data_path.exists():
            raw = json.loads(data_path.read_text(encoding="utf-8"))
            headers = raw.get("headers", list(raw.keys()))
            rows = raw.get("rows", [raw.get("data", [])])
            pdf.set_font("Helvetica", "B", 10)
            for h in headers:
                pdf.cell(40, 10, str(h), border=1)
            pdf.ln()
            pdf.set_font("Helvetica", size=10)
            for row in rows:
                for cell in row:
                    pdf.cell(40, 10, str(cell), border=1)
                pdf.ln()

    pdf.output(args.output)
    print(f"PDF generated: {args.output}")


def cmd_reportlab(args):
    """Generate a PDF using ReportLab."""
    from reportlab.lib.pagesizes import A4, letter
    from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
    from reportlab.lib.units import mm
    from reportlab.platypus import (
        SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle, PageBreak
    )
    from reportlab.lib import colors

    page_size = A4 if args.pagesize == "A4" else letter
    doc = SimpleDocTemplate(
        args.output, pagesize=page_size,
        rightMargin=20*mm, leftMargin=20*mm,
        topMargin=20*mm, bottomMargin=20*mm,
    )
    styles = getSampleStyleSheet()
    elements = []

    if args.title:
        elements.append(Paragraph(args.title, styles["Title"]))
        elements.append(Spacer(1, 12))

    if args.content:
        content_path = Path(args.content)
        if content_path.exists():
            text = content_path.read_text(encoding="utf-8")
        else:
            text = args.content
        for paragraph in text.split("\n\n"):
            elements.append(Paragraph(paragraph.strip(), styles["Normal"]))
            elements.append(Spacer(1, 6))

    if args.data:
        data_path = Path(args.data)
        if data_path.exists():
            raw = json.loads(data_path.read_text(encoding="utf-8"))
            headers = raw.get("headers", [])
            rows = raw.get("rows", raw.get("data", []))
            table_data = [headers] + rows
            table = Table(table_data)
            table.setStyle(TableStyle([
                ("BACKGROUND", (0, 0), (-1, 0), colors.HexColor("#2F5496")),
                ("TEXTCOLOR", (0, 0), (-1, 0), colors.white),
                ("FONTNAME", (0, 0), (-1, 0), "Helvetica-Bold"),
                ("FONTSIZE", (0, 0), (-1, 0), 10),
                ("FONTSIZE", (0, 1), (-1, -1), 9),
                ("GRID", (0, 0), (-1, -1), 0.5, colors.grey),
                ("ROWBACKGROUNDS", (0, 1), (-1, -1), [colors.white, colors.HexColor("#F2F2F2")]),
                ("VALIGN", (0, 0), (-1, -1), "MIDDLE"),
                ("TOPPADDING", (0, 0), (-1, -1), 6),
                ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
            ]))
            elements.append(table)

    doc.build(elements)
    print(f"PDF generated: {args.output}")


def cmd_weasyprint(args):
    """Generate a PDF from HTML using WeasyPrint."""
    import weasyprint

    if args.html:
        html_path = Path(args.html)
        if html_path.exists():
            weasyprint.HTML(filename=str(html_path)).write_pdf(args.output)
        else:
            weasyprint.HTML(string=args.html).write_pdf(args.output)
    elif args.markdown:
        import markdown
        md_path = Path(args.markdown)
        md_text = md_path.read_text(encoding="utf-8")
        html = markdown.markdown(md_text, extensions=["tables", "fenced_code", "codehilite"])
        styled_html = f"""<!DOCTYPE html>
<html><head><meta charset="utf-8">
<style>
body {{ font-family: 'Helvetica', 'Arial', sans-serif; max-width: 800px; margin: 40px auto; line-height: 1.6; color: #333; }}
h1, h2, h3 {{ color: #2F5496; }}
code {{ background: #f4f4f4; padding: 2px 6px; border-radius: 3px; }}
pre {{ background: #f4f4f4; padding: 12px; border-radius: 5px; overflow-x: auto; }}
table {{ border-collapse: collapse; width: 100%; }}
th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
th {{ background: #2F5496; color: white; }}
img {{ max-width: 100%; }}
</style></head><body>{html}</body></html>"""
        weasyprint.HTML(string=styled_html).write_pdf(args.output)
    else:
        print("Error: Provide --html or --markdown input", file=sys.stderr)
        sys.exit(1)

    print(f"PDF generated: {args.output}")


def main():
    parser = argparse.ArgumentParser(description="Generate PDFs from various inputs")
    subparsers = parser.add_subparsers(dest="command", required=True)

    # fpdf2
    fp = subparsers.add_parser("fpdf2", help="Generate PDF using fpdf2")
    fp.add_argument("--output", "-o", required=True, help="Output PDF path")
    fp.add_argument("--title", help="Document title")
    fp.add_argument("--content", help="Text content (path or inline text)")
    fp.add_argument("--data", help="JSON data file with headers/rows")

    # reportlab
    rp = subparsers.add_parser("reportlab", help="Generate PDF using ReportLab")
    rp.add_argument("--output", "-o", required=True, help="Output PDF path")
    rp.add_argument("--title", help="Document title")
    rp.add_argument("--content", help="Text content (path or inline text)")
    rp.add_argument("--data", help="JSON data file with headers/rows")
    rp.add_argument("--pagesize", default="A4", choices=["A4", "letter"], help="Page size")

    # weasyprint
    wp = subparsers.add_parser("weasyprint", help="Generate PDF using WeasyPrint (HTML/Markdown)")
    wp.add_argument("--output", "-o", required=True, help="Output PDF path")
    wp.add_argument("--html", help="HTML input (path or inline)")
    wp.add_argument("--markdown", help="Markdown file path")

    args = parser.parse_args()

    if args.command == "fpdf2":
        cmd_fpdf2(args)
    elif args.command == "reportlab":
        cmd_reportlab(args)
    elif args.command == "weasyprint":
        cmd_weasyprint(args)


if __name__ == "__main__":
    main()
