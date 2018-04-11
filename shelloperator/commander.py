import time
import subprocess as sp


def run_command(cmd_lst, stdout_fn, stderr_fn):
    with sp.Popen(cmd_lst, stdout=sp.PIPE, stderr=sp.PIPE) as p:
        while True:
            so = p.stdout.readline()

            if so:
                stdout_fn(so)

            se = p.stderr.readline()

            if se:
                stderr_fn(se)

            if p.poll() != None:
                return p.returncode

            time.sleep(1)
