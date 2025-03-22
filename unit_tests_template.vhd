library work;
	use work.address_math.all;

entity unit_tests is
end entity;
--snippet unit_test
report "
--pasteme test_condition\
" severity note;

if not (
--pasteme test_condition\
) then 
	report "FAIL" severity error;
	exit; 
end if;

--endsnippet foreach test_condition
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