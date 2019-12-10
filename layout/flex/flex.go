package flex

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
)

// Flex lays out the children of its Layoutable based on the Data assigned to
// each child.
type Flex struct {
	target       layout.Layoutable
	rows         int
	columns      int
	hSpacing     float64
	vSpacing     float64
	hAlign       align.Alignment
	vAlign       align.Alignment
	equalColumns bool
}

// New creates a new Flex layout and sets it on the Layoutable.
func New() *Flex {
	return &Flex{
		columns:  1,
		hSpacing: layout.DefaultHSpacing,
		vSpacing: layout.DefaultVSpacing,
		hAlign:   align.Start,
		vAlign:   align.Start,
	}
}

// Columns sets the number of columns. Defaults to 1.
func (f *Flex) Columns(columns int) *Flex {
	f.columns = columns
	return f
}

// HSpacing sets the spacing between columns. Defaults to DefaultHSpacing.
func (f *Flex) HSpacing(hSpacing float64) *Flex {
	f.hSpacing = hSpacing
	return f
}

// VSpacing sets the spacing between rows. Defaults to DefaultVSpacing.
func (f *Flex) VSpacing(vSpacing float64) *Flex {
	f.vSpacing = vSpacing
	return f
}

// HAlign sets the horizontal alignment of the target within its available
// space. Defaults to Start.
func (f *Flex) HAlign(hAlign align.Alignment) *Flex {
	f.hAlign = hAlign
	return f
}

// VAlign sets the vertical alignment of the target within its available
// space. Defaults to Start.
func (f *Flex) VAlign(vAlign align.Alignment) *Flex {
	f.vAlign = vAlign
	return f
}

// EqualColumns sets whether each column should use the same amount of
// horizontal space. Defaults to false.
func (f *Flex) EqualColumns(equalColumns bool) *Flex {
	f.equalColumns = equalColumns
	return f
}

// Apply the layout to the target. A copy is made of this layout and that is
// applied to the target, so this layout may be applied to other targets.
func (f *Flex) Apply(target layout.Layoutable) {
	flex := *f
	flex.target = target
	target.SetLayout(&flex)
}

// Sizes implements the Layout interface.
func (f *Flex) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	min = f.layout(geom.Point{}, hint, false, true)
	pref = f.layout(geom.Point{}, hint, false, false)
	b := f.target.Border()
	if b != nil {
		insets := b.Insets()
		min.AddInsets(insets)
		pref.AddInsets(insets)
	}
	return min, pref, layout.MaxSize(pref)
}

// Layout implements the Layout interface.
func (f *Flex) Layout() {
	var insets geom.Insets
	b := f.target.Border()
	if b != nil {
		insets = b.Insets()
	}
	hint := f.target.FrameRect().Size
	hint.SubtractInsets(insets)
	f.layout(geom.Point{X: insets.Left, Y: insets.Top}, hint, true, false)
}

func (f *Flex) layout(location geom.Point, hint geom.Size, move, useMinimumSize bool) geom.Size {
	var totalSize geom.Size
	if f.columns > 0 {
		children := f.prepChildren(useMinimumSize)
		if len(children) > 0 {
			if f.hSpacing < 0 {
				f.hSpacing = 0
			}
			if f.vSpacing < 0 {
				f.vSpacing = 0
			}
			grid := f.buildGrid(children)
			widths := f.adjustColumnWidths(hint.Width, grid)
			f.wrap(hint.Width, grid, widths, useMinimumSize)
			heights := f.adjustRowHeights(hint.Height, grid)
			totalSize.Width += f.hSpacing * float64(f.columns-1)
			totalSize.Height += f.vSpacing * float64(f.rows-1)
			for i := 0; i < f.columns; i++ {
				totalSize.Width += widths[i]
			}
			for i := 0; i < f.rows; i++ {
				totalSize.Height += heights[i]
			}
			if move {
				if totalSize.Width < hint.Width {
					if f.hAlign == align.Middle {
						location.X += xmath.Round((hint.Width - totalSize.Width) / 2)
					} else if f.hAlign == align.End {
						location.X += hint.Width - totalSize.Width
					}
				}
				if totalSize.Height < hint.Height {
					if f.vAlign == align.Middle {
						location.Y += xmath.Round((hint.Height - totalSize.Height) / 2)
					} else if f.vAlign == align.End {
						location.Y += hint.Height - totalSize.Height
					}
				}
				f.positionChildren(location, grid, widths, heights)
			}
		}
	}
	return totalSize
}

func (f *Flex) prepChildren(useMinimumSize bool) []layout.Layoutable {
	children := f.target.ChildrenForLayout()
	for _, child := range children {
		getDataFromTarget(child).computeCacheSize(child.Sizes, geom.Size{}, useMinimumSize)
	}
	return children
}

func getDataFromTarget(target layout.Layoutable) *Data {
	if data, ok := target.LayoutData().(*Data); ok {
		return data
	}
	data := NewData()
	target.SetLayoutData(data)
	return data
}

func (f *Flex) buildGrid(children []layout.Layoutable) [][]layout.Layoutable {
	var grid [][]layout.Layoutable
	var row, column int
	f.rows = 0
	for _, child := range children {
		data := getDataFromTarget(child)
		hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
		vSpan := xmath.MaxInt(1, data.vSpan)
		for {
			lastRow := row + vSpan
			if lastRow >= len(grid) {
				grid = append(grid, make([]layout.Layoutable, f.columns))
			}
			for column < f.columns && grid[row][column] != nil {
				column++
			}
			endCount := column + hSpan
			if endCount <= f.columns {
				index := column
				for index < endCount && grid[row][index] == nil {
					index++
				}
				if index == endCount {
					break
				}
				column = index
			}
			if column+hSpan >= f.columns {
				column = 0
				row++
			}
		}
		for j := 0; j < vSpan; j++ {
			pos := row + j
			for k := 0; k < hSpan; k++ {
				grid[pos][column+k] = child
			}
		}
		f.rows = xmath.MaxInt(f.rows, row+vSpan)
		column += hSpan
	}
	return grid
}

func (f *Flex) adjustColumnWidths(width float64, grid [][]layout.Layoutable) []float64 {
	availableWidth := width - f.hSpacing*float64(f.columns-1)
	expandCount := 0
	widths := make([]float64, f.columns)
	minWidths := make([]float64, f.columns)
	expandColumn := make([]bool, f.columns)
	for j := 0; j < f.columns; j++ {
		for i := 0; i < f.rows; i++ {
			data := f.getData(grid, i, j, true)
			if data != nil {
				hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
				if hSpan == 1 {
					w := data.cacheSize.Width
					if widths[j] < w {
						widths[j] = w
					}
					if data.hGrab {
						if !expandColumn[j] {
							expandCount++
						}
						expandColumn[j] = true
					}
					minimumWidth := data.minCacheSize.Width
					if !data.hGrab {
						if minimumWidth < 1 {
							w = data.cacheSize.Width
						} else {
							w = minimumWidth
						}
						if minWidths[j] < w {
							minWidths[j] = w
						}
					}
				}
			}
		}
		for i := 0; i < f.rows; i++ {
			data := f.getData(grid, i, j, false)
			if data != nil {
				hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
				if hSpan > 1 {
					var spanWidth, spanMinWidth float64
					spanExpandCount := 0
					for k := 0; k < hSpan; k++ {
						spanWidth += widths[j-k]
						spanMinWidth += minWidths[j-k]
						if expandColumn[j-k] {
							spanExpandCount++
						}
					}
					if data.hGrab && spanExpandCount == 0 {
						expandCount++
						expandColumn[j] = true
					}
					w := data.cacheSize.Width - spanWidth - float64(hSpan-1)*f.hSpacing
					if w > 0 {
						if f.equalColumns {
							equalWidth := math.Floor((w + spanWidth) / float64(hSpan))
							for k := 0; k < hSpan; k++ {
								if widths[j-k] < equalWidth {
									widths[j-k] = equalWidth
								}
							}
						} else {
							f.apportionExtra(w, j, spanExpandCount, hSpan, expandColumn, widths)
						}
					}
					minimumWidth := data.minCacheSize.Width
					if !data.hGrab || minimumWidth != 0 {
						if !data.hGrab || minimumWidth < 1 {
							w = data.cacheSize.Width
						} else {
							w = minimumWidth
						}
						w -= spanMinWidth + float64(hSpan-1)*f.hSpacing
						if w > 0 {
							f.apportionExtra(w, j, spanExpandCount, hSpan, expandColumn, minWidths)
						}
					}
				}
			}
		}
	}
	if f.equalColumns {
		var minColumnWidth, columnWidth float64
		for i := 0; i < f.columns; i++ {
			if minColumnWidth < minWidths[i] {
				minColumnWidth = minWidths[i]
			}
			if columnWidth < widths[i] {
				columnWidth = widths[i]
			}
		}
		if width > 0 && expandCount != 0 {
			columnWidth = math.Max(minColumnWidth, math.Floor(availableWidth/float64(f.columns)))
		}
		for i := 0; i < f.columns; i++ {
			expandColumn[i] = expandCount > 0
			widths[i] = columnWidth
		}
	} else if width > 0 && expandCount > 0 {
		var totalWidth float64
		for i := 0; i < f.columns; i++ {
			totalWidth += widths[i]
		}
		c := expandCount
		for math.Abs(totalWidth-availableWidth) > 0.01 {
			delta := (availableWidth - totalWidth) / float64(c)
			for j := 0; j < f.columns; j++ {
				if expandColumn[j] {
					if widths[j]+delta > minWidths[j] {
						widths[j] += delta
					} else {
						widths[j] = minWidths[j]
						expandColumn[j] = false
						c--
					}
				}
			}
			for j := 0; j < f.columns; j++ {
				for i := 0; i < f.rows; i++ {
					data := f.getData(grid, i, j, false)
					if data != nil {
						hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
						if hSpan > 1 {
							minimumWidth := data.minCacheSize.Width
							if !data.hGrab || minimumWidth != 0 {
								var spanWidth float64
								spanExpandCount := 0
								for k := 0; k < hSpan; k++ {
									spanWidth += widths[j-k]
									if expandColumn[j-k] {
										spanExpandCount++
									}
								}
								var w float64
								if !data.hGrab || minimumWidth < 1 {
									w = data.cacheSize.Width
								} else {
									w = minimumWidth
								}
								w -= spanWidth + float64(hSpan-1)*f.hSpacing
								if w > 0 {
									f.apportionExtra(w, j, spanExpandCount, hSpan, expandColumn, widths)
								}
							}
						}
					}
				}
			}
			if c == 0 {
				break
			}
			totalWidth = 0
			for i := 0; i < f.columns; i++ {
				totalWidth += widths[i]
			}
		}
	}
	return widths
}

func (f *Flex) apportionExtra(extra float64, base, count, span int, expand []bool, values []float64) {
	if count == 0 {
		values[base] += extra
	} else {
		extraInt := int(math.Floor(extra))
		delta := extraInt / count
		remainder := extraInt - delta*count
		for i := 0; i < span; i++ {
			j := base - i
			if expand[j] {
				values[j] += float64(delta)
			}
		}
		for remainder > 0 {
			for i := 0; i < span; i++ {
				j := base - i
				if expand[j] {
					values[j]++
					remainder--
					if remainder == 0 {
						break
					}
				}
			}
		}
	}
}

func (f *Flex) getData(grid [][]layout.Layoutable, row, column int, first bool) *Data {
	target := grid[row][column]
	if target != nil {
		data := getDataFromTarget(target)
		hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
		vSpan := xmath.MaxInt(1, data.vSpan)
		var i, j int
		if first {
			i = row + vSpan - 1
			j = column + hSpan - 1
		} else {
			i = row - vSpan + 1
			j = column - hSpan + 1
		}
		if i >= 0 && i < f.rows {
			if j >= 0 && j < f.columns {
				if target == grid[i][j] {
					return data
				}
			}
		}
	}
	return nil
}

func (f *Flex) wrap(width float64, grid [][]layout.Layoutable, widths []float64, useMinimumSize bool) {
	if width > 0 {
		for j := 0; j < f.columns; j++ {
			for i := 0; i < f.rows; i++ {
				data := f.getData(grid, i, j, false)
				if data != nil {
					if data.sizeHint.Height < 1 {
						hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
						var currentWidth float64
						for k := 0; k < hSpan; k++ {
							currentWidth += widths[j-k]
						}
						currentWidth += float64(hSpan-1) * f.hSpacing
						if currentWidth != data.cacheSize.Width && data.hAlign == align.Fill || data.cacheSize.Width > currentWidth {
							data.computeCacheSize(grid[i][j].Sizes, geom.Size{Width: math.Max(data.minCacheSize.Width, currentWidth)}, useMinimumSize)
							minimumHeight := data.minSize.Height
							if data.vGrab && minimumHeight > 0 && data.cacheSize.Height < minimumHeight {
								data.cacheSize.Height = minimumHeight
							}
						}
					}
				}
			}
		}
	}
}

func (f *Flex) adjustRowHeights(height float64, grid [][]layout.Layoutable) []float64 {
	availableHeight := height - f.vSpacing*float64(f.rows-1)
	expandCount := 0
	heights := make([]float64, f.rows)
	minHeights := make([]float64, f.rows)
	expandRow := make([]bool, f.rows)
	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.columns; j++ {
			data := f.getData(grid, i, j, true)
			if data != nil {
				vSpan := xmath.MaxInt(1, xmath.MinInt(data.vSpan, f.rows))
				if vSpan == 1 {
					h := data.cacheSize.Height
					if heights[i] < h {
						heights[i] = h
					}
					if data.vGrab {
						if !expandRow[i] {
							expandCount++
						}
						expandRow[i] = true
					}
					minimumHeight := data.minSize.Height
					if !data.vGrab || minimumHeight != 0 {
						var h float64
						if !data.vGrab || minimumHeight < 1 {
							h = data.minCacheSize.Height
						} else {
							h = minimumHeight
						}
						if minHeights[i] < h {
							minHeights[i] = h
						}
					}
				}
			}
		}
		for j := 0; j < f.columns; j++ {
			data := f.getData(grid, i, j, false)
			if data != nil {
				vSpan := xmath.MaxInt(1, xmath.MinInt(data.vSpan, f.rows))
				if vSpan > 1 {
					var spanHeight, spanMinHeight float64
					spanExpandCount := 0
					for k := 0; k < vSpan; k++ {
						spanHeight += heights[i-k]
						spanMinHeight += minHeights[i-k]
						if expandRow[i-k] {
							spanExpandCount++
						}
					}
					if data.vGrab && spanExpandCount == 0 {
						expandCount++
						expandRow[i] = true
					}
					h := data.cacheSize.Height - spanHeight - float64(vSpan-1)*f.vSpacing
					if h > 0 {
						if spanExpandCount == 0 {
							heights[i] += h
						} else {
							delta := h / float64(spanExpandCount)
							for k := 0; k < vSpan; k++ {
								if expandRow[i-k] {
									heights[i-k] += delta
								}
							}
						}
					}
					minimumHeight := data.minSize.Height
					if !data.vGrab || minimumHeight != 0 {
						var h float64
						if !data.vGrab || minimumHeight < 1 {
							h = data.minCacheSize.Height
						} else {
							h = minimumHeight
						}
						h -= spanMinHeight + float64(vSpan-1)*f.vSpacing
						if h > 0 {
							f.apportionExtra(h, i, spanExpandCount, vSpan, expandRow, minHeights)
						}
					}
				}
			}
		}
	}
	if height > 0 && expandCount > 0 {
		var totalHeight float64
		for i := 0; i < f.rows; i++ {
			totalHeight += heights[i]
		}
		c := expandCount
		delta := (availableHeight - totalHeight) / float64(c)
		for math.Abs(totalHeight-availableHeight) > 0.01 {
			for i := 0; i < f.rows; i++ {
				if expandRow[i] {
					if heights[i]+delta > minHeights[i] {
						heights[i] += delta
					} else {
						heights[i] = minHeights[i]
						expandRow[i] = false
						c--
					}
				}
			}
			for i := 0; i < f.rows; i++ {
				for j := 0; j < f.columns; j++ {
					data := f.getData(grid, i, j, false)
					if data != nil {
						vSpan := xmath.MaxInt(1, xmath.MinInt(data.vSpan, f.rows))
						if vSpan > 1 {
							minimumHeight := data.minSize.Height
							if !data.vGrab || minimumHeight != 0 {
								var spanHeight float64
								spanExpandCount := 0
								for k := 0; k < vSpan; k++ {
									spanHeight += heights[i-k]
									if expandRow[i-k] {
										spanExpandCount++
									}
								}
								var h float64
								if !data.vGrab || minimumHeight < 1 {
									h = data.minCacheSize.Height
								} else {
									h = minimumHeight
								}
								h -= spanHeight + float64(vSpan-1)*f.vSpacing
								if h > 0 {
									f.apportionExtra(h, i, spanExpandCount, vSpan, expandRow, heights)
								}
							}
						}
					}
				}
			}
			if c == 0 {
				break
			}
			totalHeight = 0
			for i := 0; i < f.rows; i++ {
				totalHeight += heights[i]
			}
			delta = (availableHeight - totalHeight) / float64(c)
		}
	}
	return heights
}

func (f *Flex) positionChildren(location geom.Point, grid [][]layout.Layoutable, widths, heights []float64) {
	gridY := location.Y
	for i := 0; i < f.rows; i++ {
		gridX := location.X
		for j := 0; j < f.columns; j++ {
			data := f.getData(grid, i, j, true)
			if data != nil {
				hSpan := xmath.MaxInt(1, xmath.MinInt(data.hSpan, f.columns))
				vSpan := xmath.MaxInt(1, data.vSpan)
				var cellWidth, cellHeight float64
				for k := 0; k < hSpan; k++ {
					cellWidth += widths[j+k]
				}
				for k := 0; k < vSpan; k++ {
					cellHeight += heights[i+k]
				}
				cellWidth += f.hSpacing * float64(hSpan-1)
				childX := gridX
				childWidth := math.Min(data.cacheSize.Width, cellWidth)
				switch data.hAlign {
				case align.Middle:
					childX += math.Max(0, (cellWidth-childWidth)/2)
				case align.End:
					childX += math.Max(0, cellWidth-childWidth)
				case align.Fill:
					childWidth = cellWidth
				default:
				}
				cellHeight += f.vSpacing * float64(vSpan-1)
				childY := gridY
				childHeight := math.Min(data.cacheSize.Height, cellHeight)
				switch data.vAlign {
				case align.Middle:
					childY += math.Max(0, (cellHeight-childHeight)/2)
				case align.End:
					childY += math.Max(0, cellHeight-childHeight)
				case align.Fill:
					childHeight = cellHeight
				default:
				}
				child := grid[i][j]
				if child != nil {
					child.SetFrameRect(geom.Rect{Point: geom.Point{X: childX, Y: childY}, Size: geom.Size{Width: childWidth, Height: childHeight}})
				}
			}
			gridX += widths[j] + f.hSpacing
		}
		gridY += heights[i] + f.vSpacing
	}
}
