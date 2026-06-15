---
name: pdf-generator
description: >
  Generate, create, or convert documents to PDF format. 
  Use this skill whenever the user wants to create a PDF from text, Markdown, HTML, JSON data, CSV tables, images, or any structured content.
  Also use it when the user needs to generate invoices, reports, certificates, forms, labels, or any printable document.
  This skill handles library selection, writing generation scripts, and producing the final PDF file in the workspace.
---

# PDF Generator Skill

A general-purpose skill for generating PDF documents from various input formats. It selects the appropriate approach based on the user's needs, installs dependencies, writes generation scripts, and produces polished PDF output.

---

## Quick Start

When the user asks to generate a PDF:

1. **Understand the input**: What is the source material? (text, markdown, HTML, CSV, JSON, images, etc.)
2. **Choose the right tool** (see "Tool Selection Guide" below)
3. **Write and run a generation script** that produces the PDF
4. **Deliver the PDF** to the user at the specified output path

**Always install dependencies in a way that works with the current environment.** Prefer `pip install <package> --quiet` before writing the generation script.

---

## Tool Selection Guide

Choose the appropriate library based on the user's needs:

| Use Case | Recommended Library | Why |
|---|---|---|
| Simple text + basic formatting | **fpdf2** (`fpdf`) | Lightweight, no HTML/CSS overhead, easy API |
| Rich layout (tables, images, headers/footers) | **reportlab** | Powerful layout engine, perfect for complex multi-page docs |
| HTML/CSS to PDF (converting web content) | **weasyprint** or **pdfkit** | Best for HTML-to-PDF with CSS styling |
| Markdown to PDF | **weasyprint** (convert MD→HTML→PDF) or **markdown** + **reportlab** | Combine markdown parser with PDF renderer |
| Simple tables/data from CSV/JSON | **fpdf2** or **reportlab** | Both handle tabular data well |
| Charts + reports | **reportlab** + **matplotlib** | Embed charts alongside text |
| Image-heavy documents | **reportlab** or **img2pdf** | Preserve image quality |

### Fallback priority

1. Try `fpdf2` first for simple PDFs (fastest to set up)
2. Use `reportlab` for complex layouts
3. Use `weasyprint` when HTML/CSS source is available
4. Use `img2pdf` when merging images into a PDF

---

## Working with fpdf2

> ⚠️ **fpdf2 v2.7+ API Note**: In fpdf2 2.7+, use `XPos.LMARGIN` / `YPos.NEXT` enums (not string `"LMARGIN"` / `"NEXT_LINE"`) for cell positioning. Import them from `fpdf.enums`.

```python
from fpdf import FPDF
from fpdf.enums import XPos, YPos

pdf = FPDF()
pdf.add_page()

# Text
pdf.set_font("Helvetica", size=12)
pdf.cell(text="Hello World", new_x=XPos.LMARGIN, new_y=YPos.NEXT)

# Multi-line text
pdf.multi_cell(w=0, text="Long paragraph text...")

# Table
headers = ["Name", "Age", "City"]
data = [["Alice", "30", "NYC"], ["Bob", "25", "LA"]]
for header in headers:
    pdf.cell(w=40, h=10, text=header, border=1)
pdf.ln()
for row in data:
    for cell in row:
        pdf.cell(w=40, h=10, text=cell, border=1)
    pdf.ln()

pdf.output("output.pdf")
```

### fpdf2 Key Features

- **`pdf.set_font(family, style, size)`** — Set font (families: Helvetica, Times, Courier; style: '', 'B', 'I', 'U')
- **`pdf.cell(w, h, text, border, align, fill)`** — Single cell. For line break, use `new_x=XPos.LMARGIN, new_y=YPos.NEXT` (import from `fpdf.enums`)
- **`pdf.multi_cell(w, h, text)`** — Auto-wrapping multi-line text
- **`pdf.image(name, x, y, w, h)`** — Insert image (supports PNG, JPG)
- **`pdf.add_page()`** — New page
- **`pdf.set_auto_page_break(auto, margin)`** — Enable automatic page breaks
- **`pdf.set_fill_color(r, g, b)`** — Background color
- **`pdf.set_text_color(r, g, b)`** — Text color
- **`pdf.set_draw_color(r, g, b)`** — Border/line color
- **`pdf.ln(h)`** — Line break

---

## Working with ReportLab

```python
from reportlab.lib.pagesizes import A4, letter
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle, Image
from reportlab.lib.styles import getSampleStyleSheet
from reportlab.lib import colors

doc = SimpleDocTemplate("output.pdf", pagesize=A4)
styles = getSampleStyleSheet()
elements = []

# Title
elements.append(Paragraph("Report Title", styles["Title"]))
elements.append(Spacer(1, 12))

# Body text
elements.append(Paragraph("Some body text here.", styles["Normal"]))

# Table
data = [["Header1", "Header2"], ["Value1", "Value2"]]
table = Table(data)
table.setStyle(TableStyle([
    ("BACKGROUND", (0, 0), (-1, 0), colors.grey),
    ("TEXTCOLOR", (0, 0), (-1, 0), colors.whitesmoke),
    ("GRID", (0, 0), (-1, -1), 0.5, colors.black),
]))
elements.append(table)

doc.build(elements)
```

### ReportLab Key Concepts

- **`SimpleDocTemplate`** — Builds document with automatic page breaks
- **Flowables**: `Paragraph`, `Spacer`, `Table`, `Image`, `PageBreak`
- **`getSampleStyleSheet()`** — Provides Title, Heading1-6, Normal, Code styles
- **`TableStyle`** — Fine-grained table formatting (background, grid, alignment, font)
- **`pagesize`** — Use `A4` or `letter` from `reportlab.lib.pagesizes`

---

## Working with WeasyPrint (HTML/CSS → PDF)

```python
import weasyprint

html_content = """
<html><body>
<h1 style="color: navy;">Title</h1>
<p>Some content here.</p>
</body></html>
"""

weasyprint.HTML(string=html_content).write_pdf("output.pdf")
```

From a file:
```python
weasyprint.HTML(filename="input.html").write_pdf("output.pdf")
```

With CSS:
```python
weasyprint.HTML(string=html_content).write_pdf("output.pdf", stylesheets=["style.css"])
```

---

## Markdown → PDF Pipeline

When the user has Markdown content and wants a PDF:

```python
import markdown
import weasyprint

md_text = """
# Title
Some **bold** content.
"""

html = markdown.markdown(md_text, extensions=["tables", "fenced_code"])
styled_html = f"""
<html><body>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/water.css">
{html}
</body></html>
"""
weasyprint.HTML(string=styled_html).write_pdf("output.pdf")
```

---

## Common Patterns

### Adding Headers and Footers (fpdf2)

```python
class MyPDF(FPDF):
    def header(self):
        self.set_font("Helvetica", "B", 10)
        self.cell(0, 10, "Header Text", align="C")
        self.ln(20)

    def footer(self):
        self.set_y(-15)
        self.set_font("Helvetica", "I", 8)
        self.cell(0, 10, f"Page {self.page_no()}/{{nb}}", align="C")

pdf = MyPDF()
pdf.alias_nb_pages()  # enables {nb} placeholder
```

### Adding Chinese/Unicode Text (fpdf2)

```python
# fpdf2 supports Unicode by default with add_font for custom TTF
# For basic CJK, just set the font:
pdf.add_font("NotoSansSC", "", "NotoSansSC-Regular.ttf", uni=True)
pdf.set_font("NotoSansSC", size=12)
pdf.cell(text="你好，世界！")
```

### Generating PDF from JSON/CSV Data

Read data → build table → output PDF. Example approach:

1. Read JSON/CSV with `json` / `csv` modules
2. Define column headers and row data
3. Use fpdf2's `cell()` or reportlab's `Table()` to render
4. Style headers differently from data rows

---

## Best Practices

1. **Always set a proper page size** — A4 or letter depending on user's region
2. **Set margins** — Default 10mm (fpdf2) or use `SimpleDocTemplate(rightMargin=..., leftMargin=...)`
3. **Use consistent fonts** — Pick one font family throughout for professional look
4. **Handle long text** — Use `multi_cell()` (fpdf2) or `Paragraph` with defined width (reportlab)
5. **Auto page breaks** — ReportLab handles this automatically; with fpdf2, check `pdf.get_y()` before `add_page()`
6. **Color scheme** — Use a consistent 2-3 color palette (don't overuse colors)
7. **Test the script** — After writing it, run it to verify the PDF generates without errors

---

## Output

- Save PDFs to the user's workspace with a descriptive name
- If the user doesn't specify a name, derive it from the content (e.g., `report.pdf`, `invoice.pdf`, `certificate.pdf`)
- After generation, confirm the file exists with `ls -lh <path>` and inform the user

---

## Example Workflows

### Invoice from JSON
```
User: "Generate a PDF invoice from this JSON"
1. Parse the JSON (items, prices, totals)
2. Use fpdf2 with header/footer containing company info
3. Build a table for line items with total row
4. Output invoice.pdf
```

### Report with Charts
```
User: "Create a PDF report from this CSV with a bar chart"
1. Read CSV with pandas or csv module
2. Generate matplotlib chart, save as PNG
3. Use reportlab to layout: title, intro text, chart image, data table
4. Output report.pdf
```

### Certificate
```
User: "Generate a certificate of completion for John"
1. Use reportlab or fpdf2 with decorative borders
2. Large centered text for the title and recipient name
3. Add date and signature line
4. Output certificate.pdf
```
