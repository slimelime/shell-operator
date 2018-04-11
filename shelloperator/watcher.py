from concurrent import futures
from kubernetes import watch
from shelloperator import executor


def setup_watches(client, watches):
    with futures.ThreadPoolExecutor(max_workers=len(watches)) as executor:
        threads = []
        for to_watch in watches:
            t = executor.submit(watch_events, client, to_watch)
            threads.append(t)

        futures.wait(threads, return_when=futures.FIRST_EXCEPTION)

        print([x.exception() for x in threads])


def run_cmd(client, cmd, update_object):
    def fn(event):
        print("recieved event {}".format(event["object"]["metadata"]["name"]))

    return fn


def watch_events(client, to_watch):
    w = watch.Watch()
    with executor.Executor(concurrency=to_watch['concurrency'], run_fn=run_cmd(client, to_watch["command"], to_watch["updateObject"])) as runner:
        while True:
            for event in w.stream(_stream_command, client=client, path=to_watch['path']):
                runner(event['object']['metadata']['uid'], event)


def _stream_command(client, path, *args, **kwargs):
    (data, _s, _h) = client.call_api(path, "GET", {}, {"watch": "true"}, auth_settings=['BearerToken'], _preload_content=False)
    return data
