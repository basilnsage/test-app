import React from 'react';
import ListPost from "./ListPost";
import SubmitPost from "./SubmitPost";

export default () => {
    return (
        <div className="container">
            <h1>Create a Post</h1>
            <SubmitPost />
            <hr/>
            <h1>Posts</h1>
            <ListPost />
        </div>
    );
}