# UID FootGun

It looks `uuid.NewUUID` in `version1`.go version generates a UUID based on the time. 
So if you generate multiple UUIDs in a short time, they can end up being the same.

If you use New, NewRandom or NewString those look like they have way more randomnes in them.

