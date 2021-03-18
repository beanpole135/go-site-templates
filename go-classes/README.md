# go-classes
This is a collection of individual go structures/classes which can be dropped into projects for quick usage

Each class is a single *.go file which can be copied/pasted into the source dir of a project as needed.

### List of Classes

### TimeFileDB (Work-In-Progress)
* Filename: TimeFileDB.go
* Description: A simple file-based database (typically for JSON structures), which is organized by timestamps and compressed. You give the object a path to a directory to use, and it will create a directory structure under that where compressed log files are created, one per day.

Directory Structure:

* Base Directory (supplied)
   * Directories: [Year Number] ("2021" for example)
      * Directories: [Month Number in two-digit form] ("04" for example)
         * Files: [Day number in two-digit form].json.bz2 ("14.json.bz2" for example)

