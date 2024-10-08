import os

# Helper function to escape backslashes and format the file path correctly
def format_file_path(filepath):
    return f"`{filepath.replace('\\', '\\\\')}`"  # Escape backslashes for Markdown

# Function to get the appropriate markdown language tag for syntax highlighting
def get_language_tag(filename):
    extension_to_language = {
        '.tsx': 'tsx',
        '.ts': 'ts',
        '.jsx': 'jsx',
        '.js': 'javascript',
        '.mjs': 'javascript',
        '.css': 'css',
        '.json': 'json',
        '.go': 'golang'  # Added .go for Go files
    }
    ext = os.path.splitext(filename)[1]
    return extension_to_language.get(ext, '')

def generate_tree_structure(directory, prefix="", is_last=True):
    tree_structure = []
    basename = os.path.basename(directory)
    
    # Only append the base folder at the root level (avoid double nesting)
    if prefix == "":
        tree_structure.append(f"{basename}/")

    # Get the list of directories and files for the current directory
    dirs = []
    files = []
    for entry in os.listdir(directory):
        if os.path.isdir(os.path.join(directory, entry)):
            dirs.append(entry)
        elif os.path.isfile(os.path.join(directory, entry)):
            files.append(entry)

    # Sort directories and files for consistent output
    dirs.sort()
    files.sort()

    for idx, dirname in enumerate(dirs):
        sub_prefix = f"{prefix}│   " if not is_last else f"{prefix}    "
        connector = "└── " if idx == len(dirs) - 1 and not files else "├── "
        subdir = os.path.join(directory, dirname)
        tree_structure.append(f"{prefix}{connector}{dirname}/")
        tree_structure.extend(generate_tree_structure(subdir, sub_prefix, idx == len(dirs) - 1 and not files))

    for idx, filename in enumerate(files):
        connector = "└── " if idx == len(files) - 1 else "├── "
        tree_structure.append(f"{prefix}{connector}{filename}")

    return tree_structure

def write_tree_and_file_contents(directory, output_md_file):
    # Generate the directory tree structure
    tree_structure = generate_tree_structure(directory)

    # Define the file extensions to scan
    valid_extensions = {'.tsx', '.ts', '.jsx', '.css', '.json', '.go'}  # Added .go here

    # Write the project directory structure and file contents to the markdown file
    with open(output_md_file, 'w', encoding='utf-8') as md_file:
        # Write the directory structure first
        md_file.write("# Project Directory Structure\n\n")
        md_file.write("```\n")
        for line in tree_structure:
            md_file.write(line + "\n")
        md_file.write("```\n\n")

        # Now write the content of the valid files
        md_file.write("# File Contents\n\n")
        for root, _, files in os.walk(directory):
            for filename in files:
                if any(filename.endswith(ext) for ext in valid_extensions):
                    filepath = os.path.join(root, filename)
                    # Write the formatted file path wrapped in backticks
                    formatted_path = format_file_path(filepath)
                    md_file.write(f"## {formatted_path}\n")

                    # Get the appropriate language tag for syntax highlighting
                    language_tag = get_language_tag(filename)
                    md_file.write(f"```{language_tag}\n")

                    # Try to read and write the file content
                    try:
                        with open(filepath, 'r', encoding='utf-8') as file_content:
                            content = file_content.read()
                            md_file.write(content + "\n")
                    except UnicodeDecodeError:
                        md_file.write(f"Skipped: Unable to decode {filename}\n")

                    # Close the code block
                    md_file.write("```\n\n")

# Specify the directory containing the Python script
current_directory = os.path.dirname(os.path.abspath(__file__))

# Specify the output markdown file path
output_md_file = os.path.join(current_directory, 'project_structure.md')

# Write the project directory structure and file contents to the markdown file
write_tree_and_file_contents(current_directory, output_md_file)

print(f"Project structure and file contents written to {output_md_file}")
