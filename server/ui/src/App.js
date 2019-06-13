import React from "react";
import { BrowserRouter as Router, Route, Link } from "react-router-dom";
import axios from "axios";
import "./style.scss";

function Index() {
    return <h2>Home</h2>;
}

class MovieMiniature extends React.Component {
    render() {
        return (
            <li>Item from movie miniature componenent: {this.props.item.Title}</li>
        );
    }
}

class Movies extends React.Component {
    state = {
        movies: []
    }

    helper() {
        let content = "Test";

        return <div>{content}</div>;
    }

    componentDidMount() {
        axios.get("http://localhost:8400/api/v1/movies/tracked")
            .then(res => {
                const movies = res.data;
                console.log(res.data)
                this.setState({ movies });
            })
    }

    render() {
        return (
            <div>
                <h2>Movies</h2>
                <ul>
                    { this.state.movies.map(movie => <MovieMiniature item={movie} />) }
                </ul>
                {this.helper()}
            </div>
        );
    }
}

function About() {
    return <h2>About</h2>;
}

function Users() {
    return <h2>Users</h2>;
}

function Header() {
    return (
        <nav className="header">
            <ul>
                <li>
                    <Link to="/">Home</Link>
                </li>
                <li>
                    <Link to="/movies">Movies</Link>
                </li>
                <li>
                    <Link to="/about/">About</Link>
                </li>
                <li>
                    <Link to="/users/">Users</Link>
                </li>
            </ul>
        </nav>
    );
}

function AppRouter() {
    return (
        <Router>
            <div>
                <Header />

                <Route path="/" exact component={Index} />
                <Route path="/movies" component={Movies} />
                <Route path="/about/" component={About} />
                <Route path="/users/" component={Users} />
            </div>
        </Router>
    );
}

export default AppRouter;
