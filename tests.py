import os

os.chdir("asset")
os.system("go test")

os.chdir("..")

os.chdir("blockchain")
os.system("go test")

os.chdir("..")

os.chdir("wallet")
os.system("go test")