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