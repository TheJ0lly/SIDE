package blockchain

import (
	"bytes"
	"github.com/TheJ0lly/GoChain/hashtree"
	"github.com/TheJ0lly/GoChain/metadata"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"os"
	"testing"
)

func TestCreateNewBlockchain(t *testing.T) {
	_, err := CreateNewBlockchain("TestLoc")

	if err != nil {
		t.Error("Error creating blockchain\n")
		t.Fail()
	}
}

func TestBlockChain_AddData(t *testing.T) {
	BC, err := CreateNewBlockchain("TestLoc")

	if err != nil {
		t.Error("Error creating blockchain\n")
		t.Fail()
		return
	}

	err = BC.AddData("Source", "Dest", "Asset")

	if err != nil {
		t.Error("Error adding data to blockchain\n")
		t.Fail()
		return
	}

	CleanUp()
}

func TestBlockChain_GetLastMetaData(t *testing.T) {
	BC, err := CreateNewBlockchain("TestLoc")

	if err != nil {
		t.Error("Error creating blockchain\n")
		t.Fail()
		return
	}

	err = BC.AddData("Source", "Dest", "Asset")

	if err != nil {
		t.Error("Error adding data to blockchain\n")
		t.Fail()
		return
	}

	md := metadata.CreateNewMetaData("Source", "Dest", "Asset")

	if BC.GetLastMetaData().GetMetadataHash() != md.GetMetadataHash() {
		t.Error("Error metadata hash is not expected\n")
		t.Fail()
		return
	}

	CleanUp()
}

func TestBlockChain_HashTreeGeneration(t *testing.T) {
	BC, err := CreateNewBlockchain("TestLoc")

	if err != nil {
		t.Error("Error creating blockchain\n")
		t.Fail()
		return
	}

	err = BC.AddData("Source", "Dest", "Asset")

	if err != nil {
		t.Error("Error adding data to blockchain\n")
		t.Fail()
		return
	}

	root := BC.GetBlocks()[1].mCurrHash

	md := metadata.CreateNewMetaData("Source", "Dest", "Asset")

	hl := getMetaDataHashes([]*metadata.MetaData{md})

	ht := hashtree.GenerateTree(hl)

	if bytes.Compare(ht.RootHash[:], root) != 0 {
		t.Error("Error hash root is not expected\n")
		t.Fail()
		return
	}

	CleanUp()
}

func CleanUp() {
	files, _ := os.ReadDir("TestLoc")

	for _, file := range files {
		_ = os.Remove(osspecifics.CreatePath("TestLoc", file.Name()))
	}
}
