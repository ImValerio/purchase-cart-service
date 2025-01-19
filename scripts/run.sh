SRC_PATH=/mnt
APP_FILE="$SRC_PATH/app"

# Check if the app file exists
if [ ! -f "$APP_FILE" ]; then
  echo "executable file not found. Running build.sh to create it..."
  sh "$SRC_PATH/scripts/build.sh"

  if [ ! -f "$APP_FILE" ]; then
    echo "Build script failed to create the executable file. Exiting..."
    exit 1
  fi
fi

# Run the app file
echo "Running the app..."
"$APP_FILE"