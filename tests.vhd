--snippet snip_a
Snippet A
--endsnippet
Should paste Snippet A:
--pasteme snip_a

--pushcontext
--snippet snip_b
Snippet B
--endsnippet
Should paste Snippet A:
--pasteme snip_a

Should paste Snippet B:
--pasteme snip_b
--popcontext

Should not paste Snippet B:
--pasteme snip_b

--pushcontext
--snippet snip_c
Snippet C
--endsnippet

Should not paste Snippet B but should paste Snippet A:
--pasteme snip_a
--pasteme snip_b

Should paste Snippet C:
--pasteme snip_c
--popcontext