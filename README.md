# Tears of Guthix Crowdsourcing API
API is located at https://www.togcrowdsourcing.com/{endpoint}

There are two endpoints that this API serves:
 - https://www.togcrowdsourcing.com/worldinfo
     - This returns a JSON with the worlds for the week that have been crowdsourced.
     - For example: { {worldInfo1}, {worldInfo2}, ..., {worldInfoN} }
     - Each 'worldInfo' has the following properties:
       - {
         "world_number": x,
         "hits":x,
         "stream_order": "xxxxxx"
         }
       - world_number: The world associated with this worldInfo. 
       - hits: Number of times this world was seen with the attached stream_order
       - stream_order: Order in which the streams changed in Tears of Guthix for that world. 
         - For example, 'gggbbb' or 'ggbbgb'


 - https://www.togcrowdsourcing.com/lastreset
   - This returns a JSON with the last time a server reset was detected. 
   - This aims to detect the weekly RuneScape server resets, though it will sometimes detect other times during the week in which the servers have been reset. 
   - For now, this keeps a JS5 socket connection with the World 2 RuneScape server, and updates the JSON that is returned when the socket connection is broken (server has been reset).
   - This returns data that looks like the following:
     - {"reset_time":"2022-03-09T11:30:42.553816747Z","last_reset_time_unix":1646825442,"last_server_uptime":37749}
       - reset_time: Human readable format when the last server reset was detected. Usually should happen on Wednesday. 
       - last_reset_time_unix: Unix time for which the last server reset was detected. This is the more useful value. 
       - last_server_uptime: How long the server was up before the weekly reset, in seconds. Ideally, this number should equate to a week's worth of seconds, but because resets are detected outside of the normally scheduled resets sometimes, this number may be shorter than a week. 