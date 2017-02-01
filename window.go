package main

import (
	"errors"
	"fmt"
	gc "github.com/rthornton128/goncurses"
)

type RenderWindow interface {
	Rows() uint
	Cols() uint
	Clear()
	SetRow(rowIndex uint, format string, args ...interface{}) error
	DrawBorder()
}

type CellStyle struct {
	bgColour int16
	fgColour int16
}

type Cell struct {
	codePoint rune
	style     CellStyle
}

type Window struct {
	rows     uint
	cols     uint
	cells    [][]Cell
	startRow uint
	startCol uint
}

func NewWindow() *Window {
	return &Window{}
}

func (win *Window) Resize(viewDimension ViewDimension) {
	if win.rows == viewDimension.rows && win.cols == viewDimension.cols {
		return
	}

	win.rows = viewDimension.rows
	win.cols = viewDimension.cols

	win.cells = make([][]Cell, win.rows)

	for i := uint(0); i < win.rows; i++ {
		win.cells[i] = make([]Cell, win.cols)
	}
}

func (win *Window) SetPosition(startRow, startCol uint) {
	win.startRow = startRow
	win.startCol = startCol
}

func (win *Window) Rows() uint {
	return win.rows
}

func (win *Window) Cols() uint {
	return win.cols
}

func (win *Window) Clear() {
	for i := uint(0); i < win.rows; i++ {
		for j := uint(0); j < win.cols; j++ {
			win.cells[i][j].codePoint = ' '
		}
	}
}

func (win *Window) SetRow(rowIndex uint, format string, args ...interface{}) error {
	if rowIndex >= win.rows {
		return errors.New(fmt.Sprintf("Invalid row index: %v >= %v rows", rowIndex, win.rows))
	}

	str := fmt.Sprintf(format, args...)

	colIndex := uint(0)
	rowCells := win.cells[rowIndex]

	for _, codePoint := range str {
		rowCells[colIndex].codePoint = codePoint
		colIndex++

		if colIndex >= win.cols {
			break
		}
	}

	for colIndex < win.cols {
		rowCells[colIndex].codePoint = ' '
		colIndex++
	}

	return nil
}

func (win *Window) DrawBorder() {
	return

	if win.rows < 3 || win.cols < 3 {
		return
	}

	firstRow := win.cells[0]
	firstRow[0].codePoint = rune(gc.ACS_ULCORNER)

	for i := uint(1); i < win.cols-1; i++ {
		firstRow[i].codePoint = rune(gc.ACS_HLINE)
	}

	firstRow[win.cols-1].codePoint = rune(gc.ACS_URCORNER)

	for i := uint(1); i < win.rows-1; i++ {
		row := win.cells[i]
		row[0].codePoint = rune(gc.ACS_VLINE)
		row[win.rows-1].codePoint = rune(gc.ACS_VLINE)
	}

	lastRow := win.cells[win.rows-1]
	lastRow[0].codePoint = rune(gc.ACS_LLCORNER)

	for i := uint(1); i < win.cols-1; i++ {
		lastRow[i].codePoint = rune(gc.ACS_HLINE)
	}

	lastRow[win.cols-1].codePoint = rune(gc.ACS_LRCORNER)
}