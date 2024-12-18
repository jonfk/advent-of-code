package day7

import (
	"log"
	"strconv"
	"strings"
)

type operator int

const (
	ADD operator = iota
	MUL
	CONCAT
)

type equation struct {
	result int
	nums   []int
}

func Parse(input string) []equation {
	var equations []equation
	for _, line := range strings.Split(input, "\n") {
		if len(line) == 0 {
			continue
		}
		equationStr := strings.Split(line, ":")
		result, err := strconv.Atoi(equationStr[0])
		if err != nil {
			log.Fatalf("error parsing result from equation = %v, res = %v\n", line, equationStr)
		}

		numsStr := strings.Split(equationStr[1], " ")
		var nums []int
		for i := range numsStr {
			if numsStr[i] != "" {
				num, err := strconv.Atoi(numsStr[i])
				if err != nil {
					log.Fatalf("error parsing num from equation = %v, nums = %v\n", line, numsStr)
				}
				nums = append(nums, num)
			}
		}

		equations = append(equations, equation{result: result, nums: nums})
	}
	return equations
}

func TotalValidEquations(equations []equation) int {
	var totalSum int
	for i := range equations {
		opsPermutations := GenPermutationOps(len(equations[i].nums) - 1)

		for _, ops := range opsPermutations {
			if EvalEquation(ops, equations[i].nums) == equations[i].result {
				// log.Printf("Valid eq %v\n", equations[i].nums)
				totalSum += equations[i].result
				break
			}
		}
	}
	return totalSum
}

func EvalEquation(ops []operator, nums []int) int {
	var res int = nums[0]
	for i, op := range ops {
		switch op {
		case ADD:
			res = res + nums[i+1]
		case MUL:
			res = res * nums[i+1]
		case CONCAT:
			// TODO: implement concat mathematically and through string manipulation. Benchmark to see difference
			resStr := strconv.Itoa(res) + strconv.Itoa(nums[i+1])
			newRes, err := strconv.Atoi(resStr)
			if err != nil {
				log.Fatalf("Not possible Parse error Atoi. resStr = %v\n", resStr)
			}
			res = newRes
		}
	}
	return res
}

// A M
// AA MA AM MM
// AAA MAA AMA MMA AAM MAM AMM MMM
func GenPermutationOps(n int) [][]operator {
	var ops [][]operator
	for i := 0; i < n; i++ {
		if i == 0 {
			ops = [][]operator{{ADD}, {MUL}, {CONCAT}}
		} else {
			clonedOps := CloneOps(ops)
			clonedOps2 := CloneOps(ops)

			for i := range ops {
				ops[i] = append(ops[i], ADD)
			}
			for i := range clonedOps {
				clonedOps[i] = append(clonedOps[i], MUL)
			}
			for i := range clonedOps2 {
				clonedOps2[i] = append(clonedOps2[i], CONCAT)
			}
			ops = append(ops, clonedOps...)
			ops = append(ops, clonedOps2...)
		}
	}
	return ops
}

func CloneOps(ops [][]operator) [][]operator {
	clonedOps := make([][]operator, len(ops))
	copy(clonedOps, ops)
	for i := range clonedOps {
		clonedOps[i] = make([]operator, len(ops[i]))
		copy(clonedOps[i], ops[i])
	}
	return clonedOps
}
