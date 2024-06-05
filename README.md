# nat-sync
Sync video playback with your friends.

### Work in progress.

nat-sync is a server-client program that syncronizes media playback between clients.

It comes out of a great love of syncplay, and a desire to make it feel native and great to use.
We won't be forking or using any code from syncplay. It's written with thousands of lines of, I assume,
beautiful Python. Syncplay is one of the best user-space Python applications I've ever used, but it's
much more interesting to me to try to use Go. It's kind of made for an application like this.

The idea is to make the client-server sync working, and then make a good TUI for it for Linux.
After that, get the TUI to work on MacOS, then a GUI client for Windows, Linux, Web, and MacOS, in that order. 

I'm interested in using the "proper" languages for the GUIs, C# for Windows with their library, 
Swift for MacOS, and C/C++ GTK for Linux. 

For the TUI, I'm going to use Go, for two reasons. One, it's simpler to build it in the same language as the server itself, making a viable release happen sooner. 
Two, Go just so happens to have a ton of TUI support with libraries like bubbletea. However, putting in the work to make it beautiful will come much later.

If anyone wants to join in on the fun, please submit an issue.

Initial design considerations:

- Server-client is probably necessary for syncing to work. A neutral party who distributes truth to clients.
If someone wants for their computer to be the server, it should be easy. If someone wants to spin up their own server, it should be easy.

- Sqlite for persistent data. It is maximally self-hostable, easy to set up, braindead syntax, and as performant as this application could ever need.

- The actual messaging protocol for sending clients commands is going to be the bulk of the operation, I assume. I expect I'll make a bespoke binary or ascii protocol over TCP. 

- The only video players I personally care about are mpv and vlc. I will make a good interface for adding on others as we go along.
Getting all of the video players to work with all of the different GUIs will be a challenge. I expect to be using gRPC for talking to
the different languages, which should make this process easier. I haven't done this before, so it will be interesting to learn.

- natsync.json will configure the server. Any arguments in it will be available as command line options and vice versa. 

Warning: this will hopefully be over-engineered. It's a learning experience for me, recreating an application that I love.
Learning best practices by building and seeing what works.

While I'm not going to copy or translate any code from syncplay, I will make notes of when I was inspired from what they did.
Hopefully my MIT license and their Apache 2.0 license is compatible with this. As far as I can tell, I don't have responsibility towards their license
because I'm not at all modifying or taking their code. This is a full-scale rewrite, mostly from the outside looking in - making assumptions on how it
works, or how it should work. With this said, in any instance in which I do get specifically inspired towards a solution, I am going to attribute in a comment in the relevant file. 
Some version of this paragraph will be available in future revisions of this readme, unless I find that I did not take any such inspiration.
