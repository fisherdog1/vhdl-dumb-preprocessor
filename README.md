# VHDL Dumb Preprocessor
There are many cases where VHDL sources could do with some pre-processing. Certain code is just tedious to write, especially if language support is limited or incorrect. A tool may be missing VHDL 2019 "ifdef" support for example. Vendor sources also tend to have deviations from modern VHDL standards, which are impractical to fix by hand because they are numerous and are often in auto-generated files. . I currently envision the following use cases:
* Generating VHDL unit tests from a template and a list of expressions
* Generating various initialized structures like ROMs
* Manually expanding instantiations that use features Vivado doesn't like
* Generating testbenches

# Features That Actually Exist that You Can Use Right Now
* Copy and paste snippets from files or the command line
* Generate new snippets by parameterizing them with other snippets
* "pragma push/pop" style scoping of snippets

## Features I could be Bothered to Add
* Extraction/insertion of documentation from entities and subprograms
* Automatic checking of outputs for valid syntax, probably just by calling GHDL
* gcc style ifdef
* Parameters for snippets, making them more like gcc macros. However, you can kind of already do this with push/pop

## Example
Generate unit tests from a template and list of conditions:

template.vhd
```
--snippet unit_test
report "
--pasteme conditions\
" severity note;
if not (
--pasteme conditions\
) then 
	report "FAIL" severity error;
	exit; 
end if;

--endsnippet foreach conditions

library work;
	use work.address_math.all;

entity unit_tests is
end entity;

architecture sim of unit_tests is
begin
	process
	begin
		loop
			--pasteme unit_test+
			assert false report "Tests !PASS!" severity failure;
		end loop;

		assert false report "Tests !FAIL!" severity failure;
	end process;
end architecture sim;
```

conditions.vhd
```
--snippets conditions
my_function(1) == 2
my_function(3) == 4
--endsnippets
```

Run the following command:
```
.\vdp.exe -f conditions.vhd -f template.vhd -o
```
* First -f loads each line from conditions.vhd as its own snippet named "conditions"
* Second -f loads template.vhd, creates one "unit_test" snippet for each existing "conditions" snippet. A snippet named "template.vhd" is created implicitly
* -o generates code from the snippet of the last loaded file. The "pasteme" line is populated with each snippet having name "unit_test"

Output:
```
library work;
	use work.address_math.all;

entity unit_tests is
end entity;

architecture sim of unit_tests is
begin
	process
	begin
		loop
			report "my_function(1) == 2" severity note;
			
			if not (my_function(1) == 2) then 
				report "FAIL" severity error;
				exit; 
			end if;
			
			report "my_function(3) == 4" severity note;
			
			if not (my_function(3) == 4) then 
				report "FAIL" severity error;
				exit; 
			end if;
			
			assert false report "Tests !PASS!" severity failure;
		end loop;

		assert false report "Tests !FAIL!" severity failure;
	end process;
end architecture sim;
```

This technique can be used for all kinds of repetitive VHDL writing.

## Building
Install Go for your platform: https://go.dev/dl/
```
cd vdp
go build .
```

## Binaries
If any group can be trusted to find a way to run an executable, even if it's not for the specific platform they're on, it's FPGA people. I may add some release builds as a devops exercise.

## Why is this written in Go
Because as someone who values consistent, useful, and error-resistant syntax, Python causes significant damage to my sanity every time I use it. This choice is mostly an effort to learn Go and perhaps push back on the trend of using Python for devops tasks in hardware industries. Being less clumsy than Tcl is not a high enough bar. I would also have considered Rust, but Go was much easier to install. Also, needing a ton of extra steps to allow and check for null pointers is very annoying when writing parsers, which are traditionally a big mess of frequently-null pointers.

## Known Issues
* Cool looking crash if you pop more contexts that exist
* Frustrating addition of whitespace that should not be there
* Wobbly looking parser
