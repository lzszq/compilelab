package LL1

func GetIDTerm() map[string]int {
	IDTable := []string{"Variable", "Constant", "sin", "cos", "tg", "ctg", "lg", "log", "ln", "+", "-", "*", "/", "^", "(", ")", ",", ";", "=", "?", "#"}
	IDTerm := make(map[string]int)
	for i, item := range IDTable {
		IDTerm[item] = i
	}
	return IDTerm
}

func GetIDNonTerm() map[string]int {
	IDTable := []string{"S", "A", "A'", "B", "B'", "C", "C'", "C''", "D", "D'"}
	IDNonTerm := make(map[string]int)
	for i, item := range IDTable {
		IDNonTerm[item] = i
	}
	return IDNonTerm
}

func GetPredictTable() [][][]string {
	PredictTable := make([][][]string, 10)
	for i := range PredictTable {
		PredictTable[i] = make([][]string, 21)
	}
	PredictTable[0][0] = append(PredictTable[0][0], ";", "A", "=", "Variable")
	PredictTable[0][19] = append(PredictTable[0][19], ";", "A", "?")
	PredictTable[0][20] = append(PredictTable[0][20], "ε")

	PredictTable[1][0] = append(PredictTable[1][0], "A'", "B")
	PredictTable[1][1] = append(PredictTable[1][1], "A'", "B")
	PredictTable[1][2] = append(PredictTable[1][2], "A'", "B")
	PredictTable[1][3] = append(PredictTable[1][3], "A'", "B")
	PredictTable[1][4] = append(PredictTable[1][4], "A'", "B")
	PredictTable[1][5] = append(PredictTable[1][5], "A'", "B")
	PredictTable[1][6] = append(PredictTable[1][6], "A'", "B")
	PredictTable[1][7] = append(PredictTable[1][7], "A'", "B")
	PredictTable[1][8] = append(PredictTable[1][8], "A'", "B")
	PredictTable[1][9] = append(PredictTable[1][9], "A'", "B")
	PredictTable[1][10] = append(PredictTable[1][10], "A'", "B")
	PredictTable[1][14] = append(PredictTable[1][14], "A'", "B")

	PredictTable[2][9] = append(PredictTable[2][9], "A'", "B", "+")
	PredictTable[2][10] = append(PredictTable[2][10], "A'", "B", "-")
	PredictTable[2][15] = append(PredictTable[2][15], "ε")
	PredictTable[2][17] = append(PredictTable[2][17], "ε")

	PredictTable[3][0] = append(PredictTable[3][0], "B'", "C")
	PredictTable[3][1] = append(PredictTable[3][1], "B'", "C")
	PredictTable[3][2] = append(PredictTable[3][2], "B'", "C")
	PredictTable[3][3] = append(PredictTable[3][3], "B'", "C")
	PredictTable[3][4] = append(PredictTable[3][4], "B'", "C")
	PredictTable[3][5] = append(PredictTable[3][5], "B'", "C")
	PredictTable[3][6] = append(PredictTable[3][6], "B'", "C")
	PredictTable[3][7] = append(PredictTable[3][7], "B'", "C")
	PredictTable[3][8] = append(PredictTable[3][8], "B'", "C")
	PredictTable[3][9] = append(PredictTable[3][9], "B'", "C")
	PredictTable[3][10] = append(PredictTable[3][10], "B'", "C")
	PredictTable[3][14] = append(PredictTable[3][14], "B'", "C")

	PredictTable[4][9] = append(PredictTable[4][9], "ε")
	PredictTable[4][10] = append(PredictTable[4][10], "ε")
	PredictTable[4][11] = append(PredictTable[4][11], "B'", "C", "*")
	PredictTable[4][12] = append(PredictTable[4][12], "B'", "C", "/")
	PredictTable[4][15] = append(PredictTable[4][15], "ε")
	PredictTable[4][17] = append(PredictTable[4][17], "ε")

	PredictTable[5][0] = append(PredictTable[5][0], "C''", "D")
	PredictTable[5][1] = append(PredictTable[5][1], "C''", "D")
	PredictTable[5][2] = append(PredictTable[5][2], ")", "D", "(", "sin")
	PredictTable[5][3] = append(PredictTable[5][3], ")", "D", "(", "cos")
	PredictTable[5][4] = append(PredictTable[5][4], ")", "D", "(", "tg")
	PredictTable[5][5] = append(PredictTable[5][5], ")", "D", "(", "ctg")
	PredictTable[5][6] = append(PredictTable[5][6], ")", "D", "(", "lg")
	PredictTable[5][7] = append(PredictTable[5][7], "C'", "D", "(", "log")
	PredictTable[5][8] = append(PredictTable[5][8], ")", "D", "(", "ln")
	PredictTable[5][9] = append(PredictTable[5][9], "C''", "D")
	PredictTable[5][10] = append(PredictTable[5][10], "C''", "D")
	PredictTable[5][14] = append(PredictTable[5][14], "C''", "D")

	PredictTable[6][15] = append(PredictTable[6][15], ")")
	PredictTable[6][16] = append(PredictTable[6][16], ")", "D", ",")

	PredictTable[7][9] = append(PredictTable[7][9], "ε")
	PredictTable[7][10] = append(PredictTable[7][10], "ε")
	PredictTable[7][11] = append(PredictTable[7][11], "ε")
	PredictTable[7][12] = append(PredictTable[7][12], "ε")
	PredictTable[7][13] = append(PredictTable[7][13], "D", "^")
	PredictTable[7][15] = append(PredictTable[7][15], "ε")
	PredictTable[7][17] = append(PredictTable[7][17], "ε")

	PredictTable[8][0] = append(PredictTable[8][0], "D'")
	PredictTable[8][1] = append(PredictTable[8][1], "D'")
	PredictTable[8][9] = append(PredictTable[8][9], "D'", "+")
	PredictTable[8][10] = append(PredictTable[8][10], "D'", "-")
	PredictTable[8][14] = append(PredictTable[8][14], ")", "A", "(")

	PredictTable[9][0] = append(PredictTable[9][0], "Variable")
	PredictTable[9][1] = append(PredictTable[9][1], "Constant")

	return PredictTable
}
