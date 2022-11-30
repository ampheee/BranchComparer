                                            ## BranchComparer ##
## How to use? ##
Use the following CLI commands:

`make` - to execute the target in the Makefile utility and compile the startup file


`./branchComparer` - to run the file
When the program starts, enter the branches with a space:

Example: `p10 sisyphus`

<img src="https://i.imgur.com/y5qQz18.png">

Next - just wait for the data to be processed and go to the directory of the executable file.
(You will get the appropriate message with metrics of time-usage).

<img src="https://i.imgur.com/J4UNHju.png">

In the directory itself the files will look like:

*Uniq.json - 2 files with the names of the corresponding branches, which stores the unique package branches
Updated.json - 1 file with updated packages

<img src="https://i.imgur.com/PCZZtqM.png">

!!! The proposed answer to the test task does not provide error handling for non-existent branches!!!

by Nazipov Rustam