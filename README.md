# Awesome Autostart Manager (AAM)

## Description

AAM is a GUI application designed to manage startup apps in Windows. With this application, you can easily control which apps start up when your Windows machine boots, helping to improve system performance and user experience. This is the initial Alpha version and is designed for my own usage.

## Usage

![Alt Text](https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExMmE1YjMxMzI5ZGZiNzUwZDQ0MGZmZDkyZGM5N2FkNjY4NTM5YjM5MCZlcD12MV9pbnRlcm5hbF9naWZzX2dpZklkJmN0PWc/U8iMFtS32920OmMshz/giphy.gif)

## Features

- **Autostart Application Overview**: The main window displays a table of all applications set to autostart on your Windows system. This includes both applications registered in the registry under `Software\Microsoft\Windows\CurrentVersion\Run` and shortcuts placed in the user's Startup folder.

- **Add Autostart Applications**: You can add new applications to the autostart list. The application lets you select an `.exe` file via a file dialog, which it then adds to your system's autostart settings.

## Known Limitations

Being an alpha version, this application comes with some limitations:

- The application has only been tested on a limited number of Windows configurations.
- The functionality for technically everything may not work properly.
- The GUI has been kept simple and lacks advanced user interface features.

## Future Developments

- More robust error handling and reporting.
- Broad testing on various Windows configurations.

## Usage

- You can download .exe file directly from repo or build it yourself using powershell

1. **Install Go [accordingly](https://go.dev/)**
  
2. **Clone repository and compile program:**

   ```
   git clone https://github.com/paxamans/awesomeProject
   cd awesomeProject
   go build -o executable.exe
   ```
3. **Run program if no errors ecountered:**

   ```powershell
   ./executable.exe
   ```
   
## Contributing

Contributions to AAM are welcome! Please see our [Pull Request](https://github.com/paxamans/awesomeProject/pulls) page for more details.

## Bug Reports

If you encounter a problem with the software, please file a bug report [here](https://github.com/paxamans/awesomeProject/issues).

## Credits

AAM is developed by [paxamans](https://github.com/paxamans).
