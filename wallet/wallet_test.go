package wallet

import "testing"

func TestCreateNewWallet_Correct(t *testing.T) {
	W, err := CreateNewWallet("TestUser", "TestPass", "TestAdd", "", "8080")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	if !W.ConfirmPassword("TestPass") {
		t.Errorf("Error: Password does not match!\n")
		t.Fail()
		return
	}
}

func TestCreateNewWallet_FailAddress(t *testing.T) {
	_, err := CreateNewWallet("TestUser", "TestPass", "TestAdd", "123456789", "8080")

	if err == nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

}

func TestWallet_ExportWallet(t *testing.T) {
	W, err := CreateNewWallet("TestUser", "TestPass", "TestAdd", "", "8080")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	err = W.ExportWallet()

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}
}

func TestImportWallet(t *testing.T) {
	W, err := CreateNewWallet("TestUser", "TestPass", "TestAdd", "", "8080")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	err = W.ExportWallet()

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	W, err = ImportWallet("TestUser")

	if !W.ConfirmPassword("TestPass") {
		t.Errorf("Error: Password does not match!\n")
		t.Fail()
		return
	}

}

func TestWallet_AddAsset(t *testing.T) {
	W, err := CreateNewWallet("TestUser", "TestPass", "TestAdd", "", "8080")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	a, err := W.AddAssetFromLocal("Photo", "../testassets/index.jpg")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	if !(a.GetName() == "Photo") {
		t.Errorf("Error: Asset did not match name\n")
		t.Fail()
		return
	}

	if !W.checkAssetExists("Photo") {
		t.Errorf("Error: There is no asset with the name Photo\n")
		t.Fail()
		return
	}

	_, err = W.RemoveAsset("Photo")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}
}

func TestWallet_RemoveAsset(t *testing.T) {
	W, err := CreateNewWallet("TestUser", "TestPass", "TestRem", "", "8080")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	a, err := W.AddAssetFromLocal("Photo", "../testassets/index.jpg")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	a, err = W.RemoveAsset("Photo")

	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.Fail()
		return
	}

	if !(a.GetName() == "Photo") {
		t.Errorf("Error: Asset did not match name\n")
		t.Fail()
		return
	}

}
