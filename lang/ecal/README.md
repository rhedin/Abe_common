ECAL - Event Condition Action Language
--
ECAL is a language to create rule based system which react to events provided that a defined condition holds:

Event -> Condition -> Action

Rules are defined as event sinks and have the following form:

sink "mysink" 
    "
    A comment describing the sink.
    "
    kindmatch [ "foo", a.b.bar ],
    scopematch [ "data.read", "data.write" ],
    statematch { a : 1, b : NULL },
    priority 0,
    suppresses [ "myothersink" ]
    {
      <ECAL Code>
    }

