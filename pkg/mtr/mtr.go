package mtr

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type MtrLine struct {
	Name string
	Pos  int
	IP   string
}

func New(s string) ([]string, error) {
	_, err := exec.Command("mtr", "-v", s).Output()
	if err != nil {
		return []string{}, err
	}
	fmt.Println("Waiting for MTR results...")
	out, err := exec.Command("mtr", "--raw", s).Output()
	hops := parseOutput(out)
	return hops, nil
}

func parseOutput(b []byte) []string {
	var hops []string
	var matrix []MtrLine

	raw := strings.Split(string(b), "\n")
	for _, rec := range raw {
		if len(rec) > 0 && rec[0] == 104 {
			tuple := strings.Split(rec, " ")
			pos, _ := strconv.Atoi(tuple[1])
			matrix = append(matrix, MtrLine{Name: tuple[0], Pos: pos, IP: tuple[2]})
		}
	}

	sort.Slice(matrix, func(i, j int) bool { return matrix[i].Pos < matrix[j].Pos })

	for _, tuple := range matrix {
		hops = append(hops, tuple.IP)
	}

	return hops
}
