#!/usr/bin/env python3
import os
import re
from typing import Dict, List, Set

class GoCodeAnalyzer:
    def __init__(self, root_dir: str):
        self.root_dir = root_dir
        self.packages: Dict[str, Set[str]] = {}  # package -> set of files
        self.interfaces: Dict[str, List[str]] = {}  # interface name -> list of methods
        self.structs: Dict[str, List[str]] = {}  # struct name -> list of fields
        self.methods: Dict[str, List[str]] = {}  # struct/interface name -> list of methods
        self.imports: Dict[str, Set[str]] = {}  # file -> set of imports

    def analyze_file(self, file_path: str) -> None:
        with open(file_path, 'r') as f:
            content = f.read()

        # Extract package name
        package_match = re.search(r'package\s+(\w+)', content)
        if package_match:
            package_name = package_match.group(1)
            rel_path = os.path.relpath(file_path, self.root_dir)
            if package_name not in self.packages:
                self.packages[package_name] = set()
            self.packages[package_name].add(rel_path)

        # Extract imports
        import_matches = re.findall(r'import\s*\(\s*(.*?)\s*\)', content, re.DOTALL)
        if import_matches:
            imports = set()
            for import_block in import_matches:
                for imp in re.finditer(r'"([^"]+)"', import_block):
                    imports.add(imp.group(1))
            self.imports[file_path] = imports

        # Extract interfaces
        interface_matches = re.finditer(r'type\s+(\w+)\s+interface\s*{([^}]+)}', content)
        for match in interface_matches:
            interface_name = match.group(1)
            methods = []
            method_block = match.group(2)
            for method in re.finditer(r'(\w+)\([^)]*\)[^{]*', method_block):
                methods.append(method.group(1))
            self.interfaces[interface_name] = methods

        # Extract structs and their methods
        struct_matches = re.finditer(r'type\s+(\w+)\s+struct\s*{([^}]+)}', content)
        for match in struct_matches:
            struct_name = match.group(1)
            fields = []
            field_block = match.group(2)
            for field in re.finditer(r'(\w+)\s+[^,\n]+', field_block):
                fields.append(field.group(1))
            self.structs[struct_name] = fields

        # Extract methods
        method_matches = re.finditer(r'func\s*\((\w+)\s+\*?(\w+)\)\s+(\w+)', content)
        for match in method_matches:
            receiver_name = match.group(1)
            type_name = match.group(2)
            method_name = match.group(3)
            if type_name not in self.methods:
                self.methods[type_name] = []
            self.methods[type_name].append(method_name)

    def analyze_directory(self) -> None:
        for root, _, files in os.walk(self.root_dir):
            for file in files:
                if file.endswith('.go'):
                    full_path = os.path.join(root, file)
                    self.analyze_file(full_path)

    def print_analysis(self) -> None:
        print("\n=== Project Structure Analysis ===\n")
        
        print("Directory Structure:")
        for root, dirs, files in os.walk(self.root_dir):
            level = root.replace(self.root_dir, '').count(os.sep)
            indent = ' ' * 4 * level
            print(f"{indent}{os.path.basename(root)}/")
            subindent = ' ' * 4 * (level + 1)
            for f in sorted(files):
                if f.endswith('.go'):
                    print(f"{subindent}{f}")

        print("\nPackages and Their Files:")
        for package, files in sorted(self.packages.items()):
            print(f"\n{package}:")
            for file in sorted(files):
                print(f"    {file}")

        print("\nInterfaces and Their Methods:")
        for interface, methods in sorted(self.interfaces.items()):
            print(f"\n{interface}:")
            for method in sorted(methods):
                print(f"    {method}")

        print("\nStructs and Their Methods:")
        for struct, methods in sorted(self.methods.items()):
            if struct in self.structs:
                print(f"\n{struct}:")
                print("  Fields:")
                for field in sorted(self.structs[struct]):
                    print(f"    {field}")
                print("  Methods:")
                for method in sorted(methods):
                    print(f"    {method}")

        print("\nPackage Dependencies:")
        for file, imports in sorted(self.imports.items()):
            if imports:
                print(f"\n{os.path.relpath(file, self.root_dir)}:")
                for imp in sorted(imports):
                    print(f"    {imp}")

def main():
    project_root = "/Users/damirmukimov/projects/mcp/mini"
    analyzer = GoCodeAnalyzer(project_root)
    analyzer.analyze_directory()
    analyzer.print_analysis()

if __name__ == "__main__":
    main()
