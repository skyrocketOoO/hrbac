import os
import case

def scan_directory(directory):
    """
    Recursively scan a directory for Go files and check enums.
    """
    for root, _, files in os.walk(directory):
      if 'docs' in root:
        continue
      if 'cmd' in root:
        for file in files:
          if file.endswith(".go"):
            file_path = os.path.join(root, file)
            if case.UseEnumStartAtZero(file_path):
              print(f"Issue found in {file_path}.")

            if case.UseNewMutex(file_path):
              print(f"Issue found in {file_path} regarding Mutex initialization.")

            if case.UsePanic(file_path):
              print(f"Issue found in {file_path} regarding panic usage.")
      else:
        for file in files:
          if file.endswith(".go"):
            file_path = os.path.join(root, file)
            if case.UseInitFunc(file_path):
              print(f"Issue: init() function is missing in {file_path}.")
              
            if case.UseFmtPrint(file_path):
              print(f"Issue: fmt.Print function is present in {file_path}.")
              
              
# Main entry point
if __name__ == "__main__":
  directory_to_scan = "./"  # Change this to the directory you want to scan
  scan_directory(directory_to_scan)
