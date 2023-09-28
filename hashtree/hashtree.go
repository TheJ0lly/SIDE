package hashtree

import "crypto/sha256"

type Pair struct {
	LHash [32]byte
	RHash [32]byte
}

type Node struct {
	Hash      [32]byte
	Direction bool
}

type Tree struct {
	TreeMatrix [][][32]byte
}

func (p *Pair) generateHash() [32]byte {
	var allBytes []byte
	allBytes = append(allBytes, p.LHash[:]...)
	allBytes = append(allBytes, p.RHash[:]...)
	return sha256.Sum256(allBytes)
}

func GenerateTree(l [][32]byte, t *Tree) [32]byte {
	var newList [][32]byte

	if len(l)%2 == 1 {
		l = append(l, l[len(l)-1])
	}

	for i := 0; i < len(l); i += 2 {
		p := Pair{LHash: l[i], RHash: l[i+1]}

		newList = append(newList, p.generateHash())
	}

	if len(newList) == 1 {
		t.TreeMatrix = append(t.TreeMatrix, newList)
		return newList[0]
	} else if len(newList)%2 == 1 {
		newList = append(newList, newList[len(newList)-1])
	}
	t.TreeMatrix = append(t.TreeMatrix, newList)

	return GenerateTree(newList, t)
}

func getListOfHashesToValidate(index int, t *Tree) []Node {
	var newList []Node

	level := 0

	for {
		if index%2 == 0 {
			newList = append(newList, Node{Hash: t.TreeMatrix[level][index+1], Direction: true})
		} else {
			newList = append(newList, Node{Hash: t.TreeMatrix[level][index-1], Direction: false})
		}
		level++

		if len(t.TreeMatrix[level]) == 1 {
			break
		}

		index = index / 2
	}

	return newList
}

func ValidateData(name string, t *Tree, rootHash [32]byte) bool {
	nameHash := sha256.Sum256([]byte(name))

	var l []Node

	for i, k := range t.TreeMatrix[0] {
		if k == nameHash {
			l = getListOfHashesToValidate(i, t)
		}
	}

	var p *Pair

	for _, k := range l {
		if k.Direction {
			p = &Pair{LHash: nameHash, RHash: k.Hash}
		} else {
			p = &Pair{LHash: k.Hash, RHash: nameHash}
		}
		nameHash = p.generateHash()
	}

	return nameHash == rootHash
}

func (t *Tree) Clear() {
	t = &Tree{}
}
