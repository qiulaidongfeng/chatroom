import redis
import threading


# Channel 管理多个聊天室
class Channel:
    __pool = redis.ConnectionPool(host="localhost", port=6379, decode_responses=True)
    __r = redis.Redis(connection_pool=__pool)
    # 保存所有聊天室的从redis得到的pubsub
    __allRoom = {}
    __lock = threading.Lock()
    # 保存所有聊天室的消息
    __roomMsg = {}
    # 用于测试
    __wait = threading.Lock()

    # 创建一个聊天室
    def CreateRoom(self, name: str) -> bool:
        pubsub = None
        with self.__lock:
            tmp = self.__allRoom.get(name)
            if tmp != None:
                return False

            pubsub = self.__r.pubsub()
            self.__allRoom[name] = pubsub

        pubsub.subscribe(name)

        self.__lock.acquire()
        t = threading.Thread(target=self.getmsg_loop, args=(name, pubsub))
        t.start()
        # 利用互斥锁同步，确保线程开始运行再返回
        self.__lock.acquire()
        self.__lock.release()
        return True

    def getmsg_loop(self, name: str, pubsub):
        self.__lock.release()
        # 监听消息
        for message in pubsub.listen():
            if message["type"] == "message":
                self.__wait.acquire(blocking=False)
                self.__lock.acquire()
                if self.__roomMsg.get(name) == None:
                    self.__roomMsg[name] = [message["data"]]
                else:
                    l = self.__roomMsg[name]
                    l.append(message["data"])
                self.__lock.release()
            elif message["type"] == "unsubscribe":
                return

    # 获取聊天室的历史消息
    def GetHistory(self, name: str):
        self.__lock.acquire()
        ret = self.__roomMsg.get(name)
        self.__lock.release()
        return ret

    # 发送一条消息到聊天室
    def SendMessage(self, room: str, msg: str):
        i = self.__r.publish(room, msg)

    # 退出聊天室
    def ExitRoom(self, name: str):
        pubsub = None
        with self.__lock:
            pubsub = self.__allRoom.get(name)
            if pubsub == None:
                return
            del self.__allRoom[name]
            if name in self.__roomMsg:
                del self.__roomMsg[name]
        pubsub.unsubscribe(name)

    def waitMessage(self):
        while self.__wait.locked() == False:
            pass
        self.__wait.release()
