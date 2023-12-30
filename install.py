import platform
import os
import sys

os_name = platform.system()

# Building the GoChain main exe and the GoChain_Install & GoChainUninstall
if os_name == "Windows":
    IsGoPresent = os.getenv("GOPATH")

    if IsGoPresent is None:
        print("Go toolchain is not installed!\nPlease install the Go toolchain to proceed!\n")
        sys.exit(1)

    location = None
    if len(sys.argv) != 2:
        print("No folder has been passed as argument...")
        print("Getting AppData\\Local folder")
        AppDataLocal = os.getenv("LOCALAPPDATA")

        if AppDataLocal is None:
            print("Cannot get the AppData\\Local folder!")
            sys.exit(2)

        location = AppDataLocal + "\\GoChain"
    else:
        location = sys.argv[1]

    print("Installing GoChain_Installer...")
    err = os.system(f"go build -o {location}\\GoChain_Installer.exe installer.go")

    if err == 127:
        print("Go toolchain not installed on this machine! Consider installing the Go toolchain before installing!")
        sys.exit(2)

    print("Installing GoChain_Uninstaller...")
    os.system(f"go build -o {location}\\GoChain_Uninstaller.exe uninstaller_windows.go")

    print("Installing GoChain...")
    os.system(f"go build -o {location}\\GoChain.exe main.go")

else:
    location = None
    if len(sys.argv) != 2:
        location = os.getcwd()
        print("No folder has been passed as argument...")
        print(f"Installing everything in {location}")
    else:
        location = sys.argv[1]

    print("Installing GoChain_Installer...")
    os.system(f"go build -o {location}\\GoChain_Installer.exe installer.go")

    print("Installing GoChain_Uninstaller...")
    os.system(f"go build -o {location}\\GoChain_Uninstaller.exe uninstaller_windows.go")

    print("Installing GoChain...")
    os.system(f"go build -o {location}\\GoChain.exe main.go")


print("\nSuccessfully installed GoChain.")
print("First use \"GoChain_Installer\" to ensure a correct functionality of the GoChain.")
print("To uninstall use \"GoChain_Uninstaller\".")
print("First use \"GoChain -h\" to see what is available after using \"GoChain_Installer\".")
print("\nTODO: PUT THE FOLDER IN THE ENVIRONMENT VARIABLES TO USE THE EXECUTABLES GLOBALLY")
