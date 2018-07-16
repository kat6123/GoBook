// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
)

func main() {
	img := constructImage()
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func constructImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z))
		}
	}

	return img
}

func constructImageParallelY() *image.RGBA {
	chy := make(chan struct{})

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		go func(py int) {
			y := float64(py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				// Image point (px, py) represents complex value z.
				img.Set(px, py, mandelbrot(z))
			}
			chy <- struct{}{}
		}(py)

	}

	for py := 0; py < height; py++ {
		<-chy
	}

	return img
}

func constructImageBuffParallelY() *image.RGBA {
	chy := make(chan struct{}, height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		go func(py int) {
			y := float64(py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				// Image point (px, py) represents complex value z.
				img.Set(px, py, mandelbrot(z))
			}
			chy <- struct{}{}
		}(py)

	}

	for py := 0; py < height; py++ {
		<-chy
	}

	return img
}

func constructImageParallelX() *image.RGBA {
	chx := make(chan struct{})

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			go func() {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				// Image point (px, py) represents complex value z.
				img.Set(px, py, mandelbrot(z))

				chx <- struct{}{}
			}()
		}
	}

	for px := 0; px < width; px++ {
		<-chx
	}

	return img
}

func constructImageParallelYAndX() *image.RGBA {
	chx := make(chan struct{})
	chy := make(chan struct{})

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		go func(py int) {
			y := float64(py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {
				go func() {
					x := float64(px)/width*(xmax-xmin) + xmin
					z := complex(x, y)
					// Image point (px, py) represents complex value z.
					img.Set(px, py, mandelbrot(z))

					chx <- struct{}{}
				}()
			}
			for px := 0; px < width; px++ {
				<-chx
			}
			chy <- struct{}{}
		}(py)

	}

	for py := 0; py < height; py++ {
		<-chy
	}

	return img
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
