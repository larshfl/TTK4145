#bash script
#Makes a directory on Desktop
#Checks if an old version is present.
#Deletes the old version 
#clones the elevator project folder for group1
#runs one elevator with a hardcoded port


DIRECTORY2="../Desktop/Gruppe1"
if [ -d "$DIRECTORY2" ]; then
    echo "Folder Gruppe1 exists"
else
    echo "Making folder Gruppe1"
    mkdir Gruppe1
fi


DIRECTORY="../Desktop/Gruppe1/project-group-1"
if [ -d "$DIRECTORY" ]; then
    echo "Repo exist - deleting project-group-1"
    rm -rf  ../Desktop/Gruppe1/project-group-1
    echo "Repo project-group-1 deleted"
    #sleep 5
    echo "Cloning project-group-1"
    git -C ../Desktop/Gruppe1 clone https://github.com/TTK4145-students-2019/project-group-1.git
else
    echo "Repo don't exist - cloning project-group-1"
    git -C ../Desktop/Gruppe1 clone https://github.com/TTK4145-students-2019/project-group-1.git
fi

# echo "Starting one elevator"
# gnome-terminal --geometry=49x14+0-1080 -e 'sh -c "cd Gruppe1/project-group-1 && go run main.go -port=10001 ; cd -"'