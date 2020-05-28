import pbs



e = pbs.event()



if e.type == pbs.EXECJOB_PROLOGUE:

        pbs.logmsg(pbs.LOG_DEBUG, "event is %s" % ("EXECJOB_PROLOGUE"))

elif e.type == pbs.EXECJOB_ATTACH:

        pbs.logmsg(pbs.LOG_DEBUG, "event is %s" % ("EXECJOB_ATTACH"))

else:

        pbs.logmsg(pbs.LOG_DEBUG, "event is %s" % ("UNKNOWN"))

