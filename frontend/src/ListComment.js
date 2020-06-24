import React from 'react';

export default ({comments}) => {
    const styledComments = Object.values(comments).map(comment => {
        let content;
        if (comment.status === "approved") {
            content = comment.body;
        }
        if (comment.status === "pending") {
            content = "comment is pending approval";
        }
        if (comment.status === "rejected") {
            content = "comment has been rejected";
        }
        return <li key={comment.id}>{content}</li>;
    });
    return <ul>
        {styledComments}
    </ul>
};
