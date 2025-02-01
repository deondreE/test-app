type User = {
  id: string,
  username: string
}

type Message = {
  id: string,
  msg: string
};

type Channel = {
  id: string,
  name: string,
  messages: Message[]
}

type Server = {
  id: string,
  users: User[],
  icon: string | undefined,
  channels: Channel[]
};
