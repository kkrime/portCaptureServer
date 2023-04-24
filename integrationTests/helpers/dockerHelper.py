import subprocess
import time
import re
# Docker helpers
def startDocker():

    if allContainersRunning():
        return True

    print("Starting docker-compose")
    out = subprocess.run(["docker-compose", "-f", "../docker-compose.yml", "up", "-d"], stdout=subprocess.PIPE)

    # print(out.returncode)
    # print("error starting up docker-compose: {}".format(out.stderr))

    if out.returncode != 0:
        print("error starting up docker-compose: {}".format(out.stderr))
        exit(-1)

    start_time = time.time()

    containersRunning = allContainersRunning()

    while allContainersRunning == False:
        print("Waiting for docker-compose to start all containers")

        time.sleep(1)
        end_time = time.time()
        elapsed_time = end_time - start_time

        if elapsed_time > 30:
            print("unable to start docker")
            exit(-1)

        containersRunning = allContainersRunning()

def allContainersRunning():
        out = subprocess.run(["docker", "container", "ps"], stdout=subprocess.PIPE)
        print(out.stdout)

        port_capture_server_translation_running = re.search(".*port_capture_server_translation.*", str(out.stdout))
        port_capture_server_running = re.search(".*port_capture_server .*", str(out.stdout))

        if port_capture_server_running == True and port_capture_server_running == True:
            return True
        return False
