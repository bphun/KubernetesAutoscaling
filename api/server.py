from main import app
# from gevent import monkey
# monkey.patch_all() # we need to patch very early

if __name__ == "__main__":
    app.run()

# import cherrypy

# if __name__ == '__main__':

#     # Mount the application
#     cherrypy.tree.graft(app, "/")

#     # Unsubscribe the default server
#     cherrypy.server.unsubscribe()

#     # Instantiate a new server object
#     server = cherrypy._cpserver.Server()

#     # Configure the server object
#     server.socket_host = "0.0.0.0"
#     server.socket_port = 5000
#     server.thread_pool = 1000

#     # For SSL Support
#     # server.ssl_module            = 'pyopenssl'
#     # server.ssl_certificate       = 'ssl/certificate.crt'
#     # server.ssl_private_key       = 'ssl/private.key'
#     # server.ssl_certificate_chain = 'ssl/bundle.crt'

#     # Subscribe this server
#     server.subscribe()

#     # Start the server engine (Option 1 *and* 2)

#     cherrypy.engine.start()
#     cherrypy.engine.block()