import platform
import os
import sys

os_name = platform.system()

# Building the SIDE main exe and the SIDE_Install & SIDEUninstall
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

        location = AppDataLocal + "\\SIDE"
    else:
        location = sys.argv[1]

    print("Installing SIDE_Installer...")
    err = os.system(f"go build -ldflags=\"-s -w\" -o {location}\\SIDE_Installer.exe installer.go")

    if err != 0:
        print("Go toolchain not installed on this machine! Consider installing the Go toolchain before installing!")
        sys.exit(2)

    print("Installing SIDE_Uninstaller...")
    os.system(f"go build -ldflags=\"-s -w\" -o {location}\\SIDE_Uninstaller.exe uninstaller.go")

    print("Installing SIDE_Service...")
    os.system(f"go build -ldflags=\"-s -w\" -o {location}\\SIDE_Service.exe service.go")

    print("Installing SIDE...")
    os.system(f"go build -ldflags=\"-s -w\" -o {location}\\SIDE.exe main.go")

else:
    location = None
    if len(sys.argv) != 2:
        location = os.getcwd()
        print("No folder has been passed as argument...")
        print(f"Installing everything in {location}")
    else:
        location = sys.argv[1]

    print("Installing SIDE_Installer...")
    err = os.system(f"go build -ldflags=\"-s -w\" -o {location}/SIDE_Installer installer.go")

    if err != 0:
        print("Go toolchain not installed on this machine! Consider installing the Go toolchain before installing!")
        sys.exit(2)

    print("Installing SIDE_Uninstaller...")
    os.system(f"go build -ldflags=\"-s -w\" -o {location}/SIDE_Uninstaller uninstaller.go")

    print("Installing SIDE_Service...")
    os.system(f"go build -ldflags=\"-s -w\" -o {location}/SIDE_Service service.go")

    print("Installing SIDE...")
    os.system(f"go build -ldflags=\"-s -w\" -o {location}/SIDE main.go")


print("\nSuccessfully installed SIDE.")
print("First use \"SIDE_Installer\" to ensure a correct functionality of the SIDE.")
print("To uninstall use \"SIDE_Uninstaller\".")
print("First use \"SIDE -h\" to see what is available after using \"SIDE_Installer\".")
print("\nTODO: PUT THE FOLDER IN THE ENVIRONMENT VARIABLES TO USE THE EXECUTABLES GLOBALLY")
