import React from 'react';
import _ from 'lodash';
import Traits from 'foliant';
import {LoginButton} from 'login';
import {Component, Refs, Values} from 'data';
import {renderTo, Spinner} from 'utils';

@renderTo('[is=foliantComponent]')
class Words extends Component {
  getState() {
    return {
      defaultLang: Values.defaultLang(),
      defaultNames: Values.defaultNames(),
      defaultWords: Values.defaultWords(),
      names: Values.names(),
      words: Values.words()
    };
  }

  isReady() {
    return !_.any(this.state, _.isEmpty);
  }

  render() {
    if (!this.isReady()) return (
      <Spinner size='large' style={{minHeight: '3em'}} />
    );

    return (
      <div>
        <div className='sf-tabset'>
          <div className='sf-tabset-head'>
            <label className='active'><h3>Names</h3></label>
            <label><h3>Words</h3></label>
          </div>

          <div className='sf-tabset-body'>
            <div className='sf-tab active' style={{padding: 0}}>
              <WordsTab title='Names' lang={this.state.defaultLang}
                        default={this.state.defaultNames} current={this.state.names} />
            </div>
            <div className='sf-tab' style={{padding: 0}}>
              <WordsTab title='Words' lang={this.state.defaultLang}
                        default={this.state.defaultWords} current={this.state.words} />
            </div>
          </div>
        </div>

        <br />

        <LoginButton />
      </div>
    );
  }
}

const limit = 12;

class WordsTab extends React.Component {
  state = {
    results: [],
    error: null,
    depleted: false
  };

  generator = null;
  lastCurrent = null;

  componentWillReceiveProps(props) {
    this.props = props;
    this.refresh();
  }

  componentWillMount() {
    this.refresh();
  }

  render() {return (
    <div className='layout-row app-words'>
      {/* Left column: source words */}
      <div className='flex-1 pad space-out'>
        <h3>Source {this.props.title}</h3>
        <form onSubmit={::this.add} className='sf-label-row sf-label-dense'
              data-sf-tooltip={this.state.error} data-sf-trigger='focus'>
          <input name='word' autofocus className={`flex-11 theme-text-primary ${this.textStyle}`} />
          <button className='flex-1 fa fa-plus theme-primary' tabIndex='-1'></button>
        </form>
        <div className={`grid-4 narrow ${this.textStyle}`}>
          {_.map(this.props.current, (word, key) => (
            <SourceWord text={word} handler={() => this.drop(key)} key={key} />
          ))}
        </div>
      </div>

      {/* Right column: generated results */}
      <div className='flex-1 pad space-out'>
        <h3>Generated {this.props.title}</h3>
        <form onSubmit={::this.generate} className='sf-label-row sf-label-dense'>
          <button className='flex-1 theme-accent fa fa-refresh' tabIndex='-1'></button>
          <button className='flex-11 theme-accent layout-row layout-center'>Generate</button>
        </form>
        <div className={`grid narrow ${this.textStyle}`}>
          {_.map(this.state.results, word => (
            <GeneratedWord text={word} handler={() => this.pick(word)} key={word} />
          ))}
          {this.state.depleted ?
            <GeneratedWord text='(depleted)'/> : null}
        </div>
      </div>
    </div>
  )}

  refresh() {
    if (_.isEqual(this.props.current, this.lastCurrent)) return;
    this.generator = getGenerator(this.props.lang, this.props.current, limit);
    this.state.results = this.generator();
    this.state.depleted = this.state.results.length < limit;
    this.lastCurrent = this.props.current;
  }

  generate(event) {
    event.preventDefault();
    let results = this.generator();
    this.setState({
      results: results,
      depleted: results.length < limit
    });
  }

  // Adds the word or displays an error message.
  add(event) {
    event.preventDefault();
    let word = event.target.word.value;
    word = word.toLowerCase();

    if (!word) {
      this.setState({error: 'Please input a word'});
      return;
    }

    if (word.length < 2) {
      this.setState({error: 'The word is too short'});
      return;
    }

    if (_.contains(this.props.current, word)) {
      this.setState({error: 'This word is already in the set'});
      return;
    }

    if (!isWordValid(this.props.lang, word)) {
      this.setState({error: 'Some of these characters are not allowed in a word'});
      return;
    }

    this.setState({error: null});
    this.ref.push(word);
    event.target.word.value = '';
  }

  // Adds the given word to the store, removing it from the generated results.
  // Expects the "new" set of source words to become a union of current+word to
  // avoid refreshing the generator and the generated words, because adding a
  // previously generated word to the same source set has no effect on the total
  // output.
  pick(word) {
    if (_.contains(this.props.current, word)) return;
    let oldResults = this.state.results;

    let ref = this.ref.push(word, err => {
      if (err) {
        delete this.lastCurrent[ref.key()];
        this.setState({results: oldResults});
      }
    });

    this.lastCurrent[ref.key()] = word;
    this.setState({
      results: _.without(this.state.results, word)
    });
  }

  // Removes the given word from the store and implicitly refreshes the
  // generator and the results.
  drop(key: string) {
    this.ref.child(key).remove(err => {
      if (!err) this.ref.once('value', snap => {
        if (_.isEmpty(snap.val())) this.ref.set(this.props.default);
      });
    });
  }

  get textStyle() {return this.props.title === 'Names' ? 'text-capitalise' : 'text-lowercase'}
  get ref() {return Refs[this.props.title.toLowerCase()]()}
}

class SourceWord extends React.Component {
  render() {return (
    <div className='word'>
      <span>{this.props.text}</span>
      <button className='fa fa-times sf-button-flat fade interactive'
              onClick={this.props.handler} type='button'></button>
    </div>
  )}
}

class GeneratedWord extends React.Component {
  render() {return (
    <div className='word'>
      {this.props.handler ?
      <button className='fa fa-arrow-left sf-button-flat fade interactive'
              onClick={this.props.handler} type='button'></button> : null}
      <span className='flex-1 text-center'>{this.props.text}</span>
    </div>
  )}
}

function getGenerator(lang, words, limit) {
  let traits = getTraits(lang);
  traits.examine(_.toArray(words));
  let generator = traits.generator();

  return function() {
    let results = [];
    let word;
    while ((word = generator()) && results.length < limit) {
      if (!_.contains(words, word)) results.push(word);
    }
    return results;
  };
}

function getTraits(lang) {
  let traits = new Traits();
  if (lang.knownSounds && lang.knownSounds.length) {
    traits.knownSounds = new Traits.StringSet(lang.knownSounds);
  }
  if (lang.knownVowels && lang.knownVowels.length) {
    traits.knownVowels = new Traits.StringSet(lang.knownVowels);
  }
  return traits;
}

function isWordValid(lang, word) {
  try {
    let traits = getTraits(lang);
    traits.examine([word]);
    return true;
  } catch (err) {
    console.warn(err);
    return false;
  }
}
