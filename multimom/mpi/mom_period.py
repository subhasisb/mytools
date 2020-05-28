import pbs
import os
import sys
import time

def print_attribs(pbs_obj):
   for a in pbs_obj.attributes:
      v = getattr(pbs_obj, a)
      if (v != None) and str(v) != "":
         pbs.logmsg(pbs.LOG_DEBUG, "%s = %s" % (a,v))


#print_attribs(pbs.server())

 
e=pbs.event()
print_attribs(e)
vn = e.vnode_list

print_attribs(vn[pbs.get_local_nodename()])
 
# Setting the resource value explicitly
x = 10
pbs.logmsg(pbs.LOG_DEBUG, "foo value is %s" % (x))
e.accept(0)

