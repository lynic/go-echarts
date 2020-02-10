package charts

import (
	"io"
)

// TODO: Grid
// GridComponentOpts is the option set for grid component.
type GridOpts struct {
	// Grid 组件离容器左侧的距离。
	// left 的值可以是像 20 这样的具体像素值，可以是像 '20%' 这样相对于容器高宽的百分比
	// 也可以是 'left', 'center', 'right'。
	// 如果 left 的值为'left', 'center', 'right'，组件会根据相应的位置自动对齐。
	Left string `json:"left,omitempty"`
	// Grid 组件离容器上侧的距离。
	// top 的值可以是像 20 这样的具体像素值，可以是像 '20%' 这样相对于容器高宽的百分比
	// 也可以是 'top', 'middle', 'bottom'。
	// 如果 top 的值为'top', 'middle', 'bottom'，组件会根据相应的位置自动对齐。
	Top string `json:"top,omitempty"`
	// Grid 组件离容器右侧的距离。
	// right 的值可以是像 20 这样的具体像素值，可以是像 '20%' 这样相对于容器高宽的百分比。
	// 默认自适应。
	Right string `json:"right,omitempty"`
	// Grid 组件离容器下侧的距离。
	// bottom 的值可以是像 20 这样的具体像素值，可以是像 '20%' 这样相对于容器高宽的百分比。
	// 默认自适应
	Bottom string `json:"bottom,omitempty"`
}

type Grid struct {
	RectChart
	GridOptsList []GridOpts
	gridIndex    int
	options      []globalOptser
}

// NewKLine creates a new kline chart.
func NewGrid(routers ...RouterOpts) *Grid {
	chart := new(Grid)
	chart.initBaseOpts(routers...)
	// chart.initXYOpts()
	chart.HasXYAxis = true
	chart.HasGrid = true
	chart.gridIndex = 0
	return chart
}

func (c *Grid) SetGridOpts(optsList ...GridOpts) {
	for _, opt := range optsList {
		c.GridOptsList = append(c.GridOptsList, opt)
	}
}

// SetGlobalOptions sets options for the RectChart instance.
func (c *Grid) SetGlobalOptions(options ...globalOptser) *Grid {
	// c.RectOpts.setRectGlobalOptions(options...)
	c.BaseOpts.setBaseGlobalOptions(options...)
	c.options = options
	return c
}

func (c *Grid) setXYGlobalOptions(options ...globalOptser) {
	for i := 0; i < len(options); i++ {
		option := options[i]
		switch option.(type) {
		case XAxisOpts:
			for i := range c.XAxisOptsList {
				data := c.XAxisOptsList[i].Data
				c.XAxisOptsList[i] = option.(XAxisOpts)
				c.XAxisOptsList[i].Data = data
			}
		case YAxisOpts:
			for i := range c.YAxisOptsList {
				c.YAxisOptsList[i] = option.(YAxisOpts)
			}
		}
	}
}

// RectChart 校验器
func (c *Grid) validateOpts() {
	// apply options
	// c.setXYGlobalOptions(c.options...)

	// zoom opts?
	if c.DataZoomOptsList.Len() > 0 {
		preOpt := c.DataZoomOptsList[0]
		for i := 1; i < c.gridIndex; i++ {
			if i < c.DataZoomOptsList.Len() {
				preOpt = c.DataZoomOptsList[i]
				continue
			}
			c.SetGlobalOptions(
				DataZoomOpts{XAxisIndex: []int{0, i}, Start: preOpt.Start, End: preOpt.End},
			)
		}
	}
	// grid opts ?

	// 确保 X 轴数据不会因为设置了 XAxisOpts 而被抹除
	// c.XAxisOptsList[0].Data = rc.xAxisData

	// 确保 Y 轴数标签正确显示
	for i := 0; i < len(c.YAxisOptsList); i++ {
		c.YAxisOptsList[i].AxisLabel.Show = true
	}
	c.validateAssets(c.AssetsHost)
}

// Render renders the chart and writes the output to given writers.
func (c *Grid) Render(w ...io.Writer) error {
	// add gridOpts
	// reset datazoom
	c.insertSeriesColors(c.appendColor)
	c.validateOpts()
	return renderToWriter(c, "chart", []string{}, w...)
}

// Add Add a new chart to this Grid
func (c *Grid) Add(a ...rectCharter) {
	for i := 0; i < len(a); i++ {
		// add XAxis, assume 1
		xOpt := a[i].exportXAxisOpts()[0]
		xOpt.GridIndex = c.gridIndex
		c.ExtendXAxis(xOpt)
		// c.ExtendXAxis(a[i].exportXAxisOpts()[:1]...)
		// c.XAxisOptsList[c.gridIndex].GridIndex = c.gridIndex

		// add YAxis, assume 1
		yOpt := a[i].exportYAxisOpts()[0]
		yOpt.GridIndex = c.gridIndex
		c.ExtendYAxis(yOpt)
		// c.ExtendYAxis(a[i].exportYAxisOpts()[:1]...)
		// c.YAxisOptsList[c.gridIndex].GridIndex = c.gridIndex

		// add series
		eSeries := a[i].exportSeries()
		for j := 0; j < len(eSeries); j++ {
			eSeries[j].XAxisIndex = c.gridIndex
			eSeries[j].YAxisIndex = c.gridIndex
		}
		c.Series = append(c.Series, eSeries...)
		// increase gridIndex
		c.gridIndex++
	}
}
