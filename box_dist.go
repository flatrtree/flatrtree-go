package flatrtree

import (
	"math"
)

const earthRadiusMeters = 6_371_008.8

// PlanarBoxDist returns the squared distance between the given point and rect
func PlanarBoxDist(pX, pY, minX, minY, maxX, maxY float64) float64 {
	var dX, dY float64

	if pX < minX {
		dX = minX - pX
	} else if pX <= maxX {
		dX = 0
	} else {
		dX = pX - maxX
	}

	if pY < minY {
		dY = minY - pY
	} else if pY <= maxY {
		dY = 0
	} else {
		dY = pY - maxY
	}

	return dX*dX + dY*dY
}

// GeodeticBoxDist returns the distance in meters between the point and rect
func GeodeticBoxDist(pLon, pLat, minLon, minLat, maxLon, maxLat float64) (meters float64) {
	return earthRadiusMeters * pointRectDistGeodeticRad(
		pLat*math.Pi/180, pLon*math.Pi/180,
		minLat*math.Pi/180, minLon*math.Pi/180,
		maxLat*math.Pi/180, maxLon*math.Pi/180,
	)
}

/* Copyright (c) 2016 Josh Baker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE. */

func pointRectDistGeodeticRad(φq, λq, φl, λl, φh, λh float64) float64 {
	// Algorithm from:
	// Schubert, E., Zimek, A., & Kriegel, H.-P. (2013).
	// Geodetic Distance Queries on R-Trees for Indexing Geographic Data.
	// Lecture Notes in Computer Science, 146–164.
	// doi:10.1007/978-3-642-40235-7_9
	const (
		twoΠ  = 2 * math.Pi
		halfΠ = math.Pi / 2
	)

	// Simple case, point or invalid rect
	if φl >= φh && λl >= λh {
		return distRad(φl, λl, φq, λq)
	}

	if λl <= λq && λq <= λh {
		// q is between the bounding meridians of r
		// hence, q is north, south or within r
		if φl <= φq && φq <= φh { // Inside
			return 0
		}

		if φq < φl { // South
			return φl - φq
		}

		return φq - φh // North
	}

	// determine if q is closer to the east or west edge of r to select edge for
	// tests below
	Δλe := λl - λq
	Δλw := λq - λh
	if Δλe < 0 {
		Δλe += twoΠ
	}
	if Δλw < 0 {
		Δλw += twoΠ
	}
	var Δλ float64    // distance to closest edge
	var λedge float64 // longitude of closest edge
	if Δλe <= Δλw {
		Δλ = Δλe
		λedge = λl
	} else {
		Δλ = Δλw
		λedge = λh
	}

	sinΔλ, cosΔλ := math.Sincos(Δλ)
	tanφq := math.Tan(φq)

	if Δλ >= halfΠ {
		// If Δλ > 90 degrees (1/2 pi in radians) we're in one of the corners
		// (NW/SW or NE/SE depending on the edge selected). Compare against the
		// center line to decide which case we fall into
		φmid := (φh + φl) / 2
		if tanφq >= math.Tan(φmid)*cosΔλ {
			return distRad(φq, λq, φh, λedge) // North corner
		}
		return distRad(φq, λq, φl, λedge) // South corner
	}

	if tanφq >= math.Tan(φh)*cosΔλ {
		return distRad(φq, λq, φh, λedge) // North corner
	}

	if tanφq <= math.Tan(φl)*cosΔλ {
		return distRad(φq, λq, φl, λedge) // South corner
	}

	// We're to the East or West of the rect, compute distance using cross-track
	// Note that this is a simplification of the cross track distance formula
	// valid since the track in question is a meridian.
	return math.Asin(math.Cos(φq) * sinΔλ)
}

// distance on the unit sphere computed using Haversine formula
func distRad(φa, λa, φb, λb float64) float64 {
	if φa == φb && λa == λb {
		return 0
	}

	Δφ := φa - φb
	Δλ := λa - λb
	sinΔφ := math.Sin(Δφ / 2)
	sinΔλ := math.Sin(Δλ / 2)
	cosφa := math.Cos(φa)
	cosφb := math.Cos(φb)

	return 2 * math.Asin(math.Sqrt(sinΔφ*sinΔφ+sinΔλ*sinΔλ*cosφa*cosφb))
}
