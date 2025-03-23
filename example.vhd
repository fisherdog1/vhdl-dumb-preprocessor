
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

