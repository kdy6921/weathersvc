package main

func code2msg(code int) string {
	m := map[int]string{
		0: "clear",
		1: "cloudy",
		2: "rainy",
	}
	return m[code]
}
