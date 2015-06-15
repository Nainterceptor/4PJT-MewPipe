Suggested Video algorithm explanation
=====================================

## Client needs

    On the home page of the website, the platform will suggest users some videos the most played and shared on Mewpipe.

We have done that feature and its great.

## Our purpose

In the future, if the client wants an advanced suggestion mechanism, the dev'team can add a "Taxonomy" feature.
At upload, the user adds some tags related to the video, like "Cats" or "timelapse". 
In a video, we may also create a list of videos that have the most related taxonomies, sorting them by views-count.
If we have less than _n_ videos with all of taxonomies, the list will be filled with videos 
that matches _n-1_ taxonomies, _n-2_ taxonomies and so on...