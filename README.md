# Awesome Autostart Manager (AAM)

## Description

AAM is a GUI application designed to manage startup apps in Windows. With this application, you can easily control which apps start up when your Windows machine boots, helping to improve system performance and user experience.

## Usage

![Alt Text](https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExMmE1YjMxMzI5ZGZiNzUwZDQ0MGZmZDkyZGM5N2FkNjY4NTM5YjM5MCZlcD12MV9pbnRlcm5hbF9naWZzX2dpZklkJmN0PWc/U8iMFtS32920OmMshz/giphy.gif)

# Awesome Autostart Manager(AAM) Alpha PreRelease Notes

The AAM application is a new utility for managing autostart applications on Windows systems. This is the initial Alpha version and is designed to give a sneak peek at some of the application's core functionality.

## Features

- **Autostart Application Overview**: The main window displays a table of all applications set to autostart on your Windows system. This includes both applications registered in the registry under `Software\Microsoft\Windows\CurrentVersion\Run` and shortcuts placed in the user's Startup folder.

- **Add Autostart Applications**: You can add new applications to the autostart list. The application lets you select an `.exe` file via a file dialog, which it then adds to your system's autostart settings.

- **Delete Autostart Applications**: For each application in the table, there is a delete button. Clicking this button removes the corresponding application from your system's autostart settings.

- **Rename Autostart Applications**: For each application in the table, there is an edit button. Clicking this button allows you to rename the application entry in the autostart settings.

- **Refresh Functionality**: A refresh button is provided to update the table with the latest autostart settings.

- **Error Dialogs**: Error dialogs are provided to notify you of any problems encountered while adding, deleting, or renaming autostart applications.

## Known Limitations

Being an alpha version, this application comes with some limitations:

- There is no support for sorting or filtering the list of autostart applications.
- The application has only been tested on a limited number of Windows configurations.
- The functionality for technically everything may not work properly.
- The GUI has been kept simple and lacks advanced user interface features.

## Future Developments

We are planning a number of improvements and new features for future versions. These include:

- More robust error handling and reporting.
- Sorting and filtering of autostart applications.
- Improved user interface with advanced GUI features.
- Support for managing services as well as applications.
- Broad testing on various Windows configurations.

We hope you enjoy using the alpha version of AAM and welcome any feedback you may have!


## Contributing

Contributions to AAM are welcome! Please see our [Pull Request](https://github.com/paxamans/awesomeProject/pulls) page for more details.

## Bug Reports

If you encounter a problem with the software, please file a bug report [here](https://github.com/paxamans/awesomeProject/issues).

## Credits

AAM is developed by [paxamans](https://github.com/paxamans).
