library work;
	use work.address_math.all;

entity unit_tests is
end entity;
architecture sim of unit_tests is
begin
	process
	begin
		loop
			report "UnitTestCondition1" severity note;
			
			if not (UnitTestCondition1) then 
				report "FAIL" severity error;
				exit; 
			end if;
			
			report "UnitTestCondition2" severity note;
			
			if not (UnitTestCondition2) then 
				report "FAIL" severity error;
				exit; 
			end if;
			
			assert false report "Tests !PASS!" severity failure;
		end loop;

		assert false report "Tests !FAIL!" severity failure;
	end process;
end architecture sim;

