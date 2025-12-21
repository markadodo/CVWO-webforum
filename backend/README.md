backend design:

Data models:

    user: 
    - username
    - user id
    - password hash
    - personal bio stuff?
    - created at
    - last active at

    topic:
    - topic id
    - topic title(name)
    - topic description
    - created by(ownership)
    - created at
    - tags

    post:
    - post id
    - post title
    - post description
    - post thread(topic id)
    - post status(likes, dislikes, edited?, views, popularity)
    - created by
    - created at
    - tags
    - attachments?

    comment:
    - comment id
    - comment description
    - comment status(likes, dislikes, edited?)
    - comment thread(post id)
    - parent comment id
    - created by
    - created at