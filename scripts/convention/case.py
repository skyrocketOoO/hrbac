import re


def UseEnumStartAtZero(file_path):
    """
    Check if enums in the given Go file start at 1.
    """
    with open(file_path, "r") as f:
        lines = f.readlines()

    in_const_block = False
    for line in lines:
        line = line.strip()

        # Detect the start of a const block
        if line.startswith("const ("):
            in_const_block = True
            continue

        # Detect the end of a const block
        if in_const_block and line.endswith(")"):
            in_const_block = False

        # Look for enums using iota
        if in_const_block and "iota" in line:
            match = re.search(r"(\w+)\s*=\s*iota", line)
            if match:
                first_enum = match.group(1)
                if "= iota" in line and not "iota + 1" in line:
                    print(f"Error: Enum '{first_enum}' starts at 0 in {file_path}.")
                    return True
    return False
  
def UseNewMutex(file_path):
    """
    Check if `new(sync.Mutex)` is used in the given Go file.
    """
    with open(file_path, "r") as f:
        lines = f.readlines()

    for line in lines:
        line = line.strip()

        # Check for usage of `new(sync.Mutex)`
        if "new(sync.Mutex)" in line:
            print(f"Warning: Unnecessary use of `new(sync.Mutex)` in {file_path}. Consider using the zero value.")
            return True
    return False
  
def UsePanic(file_path):
    """
    Check if `panic` is used in the given Go file.
    """
    with open(file_path, "r") as f:
        lines = f.readlines()

    for line in lines:
        line = line.strip()

        if "panic(" in line:
            print(f"Warning: don't use panic")
            return True
    return False
  
def UseInitFunc(file_path):
    """
    Check if the Go file contains an `init()` function.
    """
    with open(file_path, "r") as f:
        lines = f.readlines()

    for line in lines:
        if line.strip().startswith("func init()"):
            print(f"Found init() function in {file_path}.")
            return True
    return False
  
def UseFmtPrint(file_path):
    """
    Check if the Go file contains an `fmt.Print` function.
    """
    with open(file_path, "r") as f:
        lines = f.readlines()

    for line in lines:
        if "fmt.Print" in line:
            print(f"Found fmt.Print function in {file_path}.")
            return True
    return False
  
  