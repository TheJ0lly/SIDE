import platform
import os
import sys

os_name = platform.system()

# Building the GoChain main exe and the GoChain_Install & GoChainUninstall
if os_name == "Windows":
    os.system("go build -o bin\\GoChain_Installer.exe installer.go")
    os.system("go build -o bin\\GoChain_Uninstaller.exe uninstaller_windows.go")
    os.system("go build -o bin\\GoChain.exe main.go")


elif os_name == "Linux":
    os.system("go build -o bin/GoChain_Installer installer.go")
    os.system("go build -o bin/GoChain_Uninstaller uninstaller_linux.go")
    os.system("go build -o bin/GoChain main.go")

elif os_name == "Darwin":
    os.system("go build -o bin/GoChain_Installer installer.go")
    os.system("go build -o bin/GoChain_Uninstaller uninstaller_ios.go")
    os.system("go build -o bin/GoChain main.go")

else:
    print(f"Unsupported OS: {os_name}\n")
    sys.exit(1)

print("Successfully installed GoChain.")
print("First use \"GoChain_Installer\" to ensure a correct functionality of the GoChain.")
print("To uninstall use \"GoChain_Uninstall\".")
print("Now you can use \"GoChain\" as you please.")
print("First use \"GoChain -h\" to see what is available.")
