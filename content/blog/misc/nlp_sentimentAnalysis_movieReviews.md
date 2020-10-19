---
title: "Sentiment analysis with Python"
date: 2020-10-19T18:54:18+11:00
draft: false
toc: true
tags: [python, nlp]
---

## Setting up

Set up your [`venv`'s](https://docs.python.org/3/library/venv.html) as your wish; then install Python's [Natural Language Toolkit (NLTK)](https://www.nltk.org/) with `pip install nltk`.

We will be looking at the [`movie_reviews`](https://www.kaggle.com/nltkdata/movie-review) corpus, but we design our classes so that they can accept any corpus.;

```
$ python

>>> import nltk
>>> nltk.download('movie_reviews')
>>> nltk.download('punkt')
>>> nltk.download('wordnet')
>>> nltk.download('stopwords')
>>> nltk.download('averaged_perceptron_tagger')
```

## NLTK corpora

A text corpus (*plural* corpora) is a large body of text.

The [`nltk.corpus`](https://github.com/nltk/nltk/tree/develop/nltk/corpus) package uses a LazyLoader so you should call `.ensure_loaded()` first. You can see the deferred initialisation by comparing the the list of callable methods before and after with: `[method for method in dir(movie_reviews) if callable(getattr(movie_reviews, method))]`.

Once the object is properly loaded, we can then use `help()` to figure out how to use it...

``` py
from nltk.corpus import movie_reviews

movie_reviews.ensure_loaded()
help(movie_reviews)
```

**n.b.** There is a `.readme()` method with important information on how the class labels were determined.
This is a worthwhile read, but I won't bother putting it here.

Lets get some descriptive statistics of our corpus:

``` py
from nltk.corpus import movie_reviews
movie_reviews.ensure_loaded()

movie_reviews.categories()
# ['neg', 'pos']

len(movie_reviews.fileids())
# 2000

len(movie_reviews.words())
# 1583820

{cat: len(movie_reviews.fileids(categories=cat))
            for cat in movie_reviews.categories()}
# Files per category
# {'neg': 1000, 'pos': 1000}

{cat: len(movie_reviews.words(categories=cat))
            for cat in movie_reviews.categories()}
# Words per category
# {'neg': 751256, 'pos': 832564}
```

## A Naive Bayes classifier

For demonstration, we'll use the `nltk.NaiveBayesClassifier`.

### Abstract base class
I've created an abstract base `CorpusModel` class for training over Corpus data:

``` py
import abc

class CorpusModel(abc.ABC):

    @abc.abstractmethod
    def train(self, data, labels):
        pass

    @abc.abstractmethod
    def predict(self, data) -> int:
        pass

    def validate(self, data, labels) -> int:

        n_correct = 0
        for d,actual in zip(data, labels):
            if self.predict(d) == actual:
                n_correct += 1

        return n_correct
```

Here we require two abstract methods to be implemented, and we will inherit the `.validate()` method.

### Preparing the training data

The `nltk.NaiveBayesClassifier` has a `.train()` method which accepts a list of tuples that look like this:

``` py
# Data format requried for training
training_data = [
    ({token1: True, token2: True, ...}, label1),    # 1st Data Point
    ({token1: True, token2: True, ...}, label2),    # 2nd Data Point
]
```

The first element of the tuple is the feature dictionary, and for our particular case we will just a token-keys with all their values set to `True`.

To get to this, we will do some pre-processing first:

1. Tokenize every review (a.k.a. data-point)
2. Filter out punctuation and stop words
3. Lemmatize the tokens

The first two steps are easy, this list comprehension will tokenize the input `data_pt` and omit any tokens that are either from `string.punctuation` or are stopwords.

``` py
import string
import nltk
#...

[token for token in nltk.word_tokenize(data_pt)
    if token not in string.punctuation
    and token.lower() not in nltk.corpus.stopwords.words('english')]
```

Lemmatizing is the process of converting a word to its base form (lemma). In this particular case, we use it to 'standardize' our sentences better.

Notice that some words can have multiple lemmas, hence specifying the POS(Part-Of-Speech) tag is necessary. As an example:

``` py
import nltk
lem = nltk.stem.wordnet.WordNetLemmatizer()

print(lem.lemmatize("stripes", "n")) # 'n' = noun
# stripe

print(lem.lemmatize("stripes", "v")) # 'v' = verb
# strip
```

Part-of-speech tagging can be achieved with `nltk.tag.pos_tag()`.
With a token-list input, we will get a list of tuples, with the first element in the tuple the token and the second the tag.

As an example, here's what happens when we tag this short sentence:

``` py
import nltk
import string

data_pt = "The quickest brown fox jumps over the lazy dogs"
tokenized = [token for token in nltk.word_tokenize(data_pt)
                if token not in string.punctuation
                and token.lower() not in nltk.corpus.stopwords.words('english')]
tagged = nltk.tag.pos_tag(tokenized)
# [('quickest', 'JJS'), ('brown', 'NN'), ('fox', 'NN'), ('jumps', 'NNS'), ('lazy', 'JJ'), ('dogs', 'NNS')]
```

The POS-tags returned here do not match the required tags used in `.lemmatize()`, hence we use a simple lookup:

``` py
import nltk
import string
from collections import defaultdict

POS_TAG_MAP = defaultdict(lambda: nltk.corpus.wordnet.NOUN)
POS_TAG_MAP["J"] = nltk.corpus.wordnet.ADJ
POS_TAG_MAP["N"] = nltk.corpus.wordnet.NOUN
POS_TAG_MAP["V"] = nltk.corpus.wordnet.VERB
POS_TAG_MAP["R"] = nltk.corpus.wordnet.ADV

# ...

lem = nltk.stem.wordnet.WordNetLemmatizer()
lemmatized = [lem.lemmatize(token, POS_TAG_MAP[pos[0].upper()]) for token, pos in tagged]
# ['quick', 'brown', 'fox', 'jump', 'lazy', 'dog']
```

Put together, we can make a method that combines this processing to return the required feature vector for training.
We add an additional dict-comprehension as required by the model.

``` py
def __init__(self):
    self._lemmatizer = nltk.stem.wordnet.WordNetLemmatizer()

def _get_feature(self, datapt):
    tagged_tokens = nltk.tag.pos_tag(
        [token for token in nltk.word_tokenize(datapt)
            if token not in string.punctuation
            and token.lower() not in nltk.corpus.stopwords.words('english')]
    )
    return {self._lemmatizer.lemmatize(token, POS_TAG_MAP[pos[0].upper()]): True
                for token, pos in tagged_tokens}
```

### Training

With a method for processing our individual data points, we are ready to implement our required `.train()` method.
This method takes in two equal sized lists - `data` and `labels`, with each element in `data` having its corresponding label in `labels`.

We first iterate over the `data` array to process it with our `_get_feature()` method, this will return the feature set of our training set.
With this, we can create a labelled `train_target` object that contains our features with their corresponding labels:

``` py
def _get_feature_set(self, data):
    return [self._get_feature(pt) for pt in data]

def train(self, data, labels):
    train_target = [(feature, lab)
                        for lab, feature in zip(labels, self._get_feature_set(data))]
    self._model = nltk.NaiveBayesClassifier.train(train_target)
```

### Predicting labels

The `nltk.NaiveBayesClassifier` has a `.classify()` method which accepts a single unlabelled feature vector, so:

``` py
def predict(self, data):
    return self._model.classify(
        self._get_feature(data)
    )
```

### The final classs
``` py
import nltk
import string
from collections import defaultdict
from . import CorpusModel

POS_TAG_MAP = defaultdict(lambda: nltk.corpus.wordnet.NOUN)
POS_TAG_MAP["J"] = nltk.corpus.wordnet.ADJ
POS_TAG_MAP["N"] = nltk.corpus.wordnet.NOUN
POS_TAG_MAP["V"] = nltk.corpus.wordnet.VERB
POS_TAG_MAP["R"] = nltk.corpus.wordnet.ADV

class NTLKNaiveBayesClassifier(CorpusModel):

    def __init__(self):
        self._lemmatizer = nltk.stem.wordnet.WordNetLemmatizer()

    def train(self, data, labels):
        train_target = [(feature, lab)
                        for lab, feature in zip(labels, self._get_feature_set(data))]
        self._model = nltk.NaiveBayesClassifier.train(train_target)

    def predict(self, data):
        return self._model.classify(
            self._get_feature(data)
        )

    def _get_feature(self, datapt):
        tagged_tokens = nltk.tag.pos_tag(
            [token for token in nltk.word_tokenize(datapt)
                if token not in string.punctuation
                and token.lower() not in nltk.corpus.stopwords.words('english')]
        )
        return {self._lemmatizer.lemmatize(token, POS_TAG_MAP[pos[0].upper()]): True
                    for token, pos in tagged_tokens}

    def _get_feature_set(self, data):
        return [self._get_feature(pt) for pt in data]
```

## Setting up a pipeline

Now we need to train our model and evaluate its performance.
To do this, we will design a `Pipeline` class which accepts a `CorpusModel`-typed object, and some corpus data. With this class we will:

1. Manage model training
2. Validation

We will use the `Pipeline` class as follows:

``` py
import nltk

model = NTLKNaiveBayesClassifier()
pl = Pipeline(model, nltk.corpus.movie_reviews)
pl.cross_validate_kfold(10)
```

Looking at the above use-case, our constructor is simple:

``` py
class Pipeline(object):

    def __init__(self, model, corpus):

        self._model = model
        self._corpus = corpus
        self._corpus.ensure_loaded()
```

### Splitting the corpus into folds:

For the `fileids` list, we can separate them into `k` equal-sized chunks with this simple `for`-loop.
We iterate over the the indexes of the `fileids` list with a step-size of `len(fileids) / k` where `k` is the number of folds we want:

``` py
folds, foldlen = [], int(len(fileids) / k)
for i in range(0, len(fileids), foldlen):
    folds.append(fileids[i:i + foldlen])
```

Once we have the folds, we can define a generator which yields a list of `fileids` for testing, and a list of `fileids` for training.
Note, the list-comprehension is to flatten the list-of-lists (`folds`):

``` py
for test_idx, test_data in enumerate(folds):
    train_data = [fileid for i,x in enumerate(folds) if i != test_idx for fileid in x]
    yield train_data, test_data
```

We put this all together into a private method of our Pipeline class, with the `shuffle` kwarg to randomize the order of the fileids:

``` py
def _kfolds(self, k, shuffle=True):

    fileids = random.sample(self._corpus.fileids(), len(self._corpus.fileids()))
                if shuffle else self._corpus.fileids()
    folds = []
    foldlen = int(len(fileids) / k)
    for i in range(0, len(fileids), foldlen):
        folds.append(fileids[i:i + foldlen])

    for test_idx, test_fileids in enumerate(folds):
        train_fileids = [fileid for i,x in enumerate(folds) if i != test_idx for fileid in x]
        yield train_fileids, test_fileids
```

### Cross-validation

To convert the `fileids` to actual data, we can create a simple method:

``` py
def _get_labelled_data(self, fileids):

    return (
        [self._corpus.raw(fid) for fid in fileids],
        [self._corpus.categories(fid)[0] for fid in fileids]
    )
```

Put together, we can iterate over all of our folds to perform cross validation as follows

``` py
def cross_validate_kfold(self, k):

    results = []
    for train_fileids, test_fileids in self._kfolds(k):

        train_data, train_labels = self._get_labelled_data(train_fileids)
        test_data, test_labels = self._get_labelled_data(test_fileids)

        self._model.train(train_data, train_labels)

        n_valid = self._model.validate(test_data, test_labels)
        n_tested = len(test_labels)
        results.append((n_tested, n_valid, n_valid/n_tested))

    return results
```

The results list is a list of 3d tuples providing us the number of tested data points, the number of valid predictions, and the corresponding ratio.

## Results

The final results of 10-folds is displayed below, we have a collective average of 71.1% accuracy

``` py
# (n_tested, n_valid, accuracy)
results = [(200, 141, 0.705),
 (200, 148, 0.74),
 (200, 153, 0.765),
 (200, 128, 0.64),
 (200, 146, 0.73),
 (200, 145, 0.725),
 (200, 131, 0.655),
 (200, 139, 0.695),
 (200, 144, 0.72),
 (200, 147, 0.735)]
```

## Discussion

We've created a simple classifier pipeline that trains any model that implements our abstract `CorpusModel` class.
We can easily substitute in a different class, as well as a different corpus from `nltk.corpus`.

There is room for improvement, both in the classifier algorithm and the way in which we manage it.
In particular, the feature extraction is a slow loop that can probably be optimised.

Ontop of that, we could look into caching models as well as multi-threading our training.

## Final code

Folder Structure:

```
root
└── pySentimentsDemo
    ├─── __init__.py
    └─── simpleModel.py
```

Run with:

```
cd root/
python -m pySentimentsDemo.simpleModel
```

**`__init__.py`**

``` py
import abc
import random

class Pipeline(object):

    def __init__(self, model, corpus):

        self._model = model
        self._corpus = corpus
        self._corpus.ensure_loaded()

    def _kfolds(self, k, shuffle=True):

        fileids = random.sample(self._corpus.fileids(), len(self._corpus.fileids())) if shuffle else self._corpus.fileids()
        folds = []
        foldlen = int(len(fileids) / k)
        for i in range(0, len(fileids), foldlen):
            folds.append(fileids[i:i + foldlen])

        for test_idx, test_fileids in enumerate(folds):
            train_fileids = [fileid for i,x in enumerate(folds) if i != test_idx for fileid in x]
            yield train_fileids, test_fileids

    def _get_labelled_data(self, fileids):

        return (
            [self._corpus.raw(fid) for fid in fileids],
            [self._corpus.categories(fid)[0] for fid in fileids]
        )

    def cross_validate_kfold(self, k):

        results = []
        for train_fileids, test_fileids in self._kfolds(k):

            train_data, train_labels = self._get_labelled_data(train_fileids)
            test_data, test_labels = self._get_labelled_data(test_fileids)

            self._model.train(train_data, train_labels)

            n_valid = self._model.validate(test_data, test_labels)
            n_tested = len(test_labels)
            results.append((n_tested, n_valid, n_valid/n_tested))

        return results

class CorpusModel(abc.ABC):

    @abc.abstractmethod
    def train(self, data, labels):
        pass

    @abc.abstractmethod
    def predict(self, data) -> int:
        pass

    def validate(self, data, labels) -> int:

        n_correct = 0
        for d,actual in zip(data, labels):
            if self.predict(d) == actual:
                n_correct += 1

        return n_correct
```

**`simpleModel.py`**

``` py
import nltk
import string
from pprint import pprint
from collections import defaultdict
from . import Pipeline, CorpusModel

POS_TAG_MAP = defaultdict(lambda: nltk.corpus.wordnet.NOUN)
POS_TAG_MAP["J"] = nltk.corpus.wordnet.ADJ
POS_TAG_MAP["N"] = nltk.corpus.wordnet.NOUN
POS_TAG_MAP["V"] = nltk.corpus.wordnet.VERB
POS_TAG_MAP["R"] = nltk.corpus.wordnet.ADV

class NTLKNaiveBayesClassifier(CorpusModel):

    def __init__(self):
        self._lemmatizer = nltk.stem.wordnet.WordNetLemmatizer()

    def train(self, data, labels):
        train_target = [(feature, lab) for lab, feature in zip(labels, self._get_feature_set(data))]
        self._model = nltk.NaiveBayesClassifier.train(train_target)

    def predict(self, data):
        return self._model.classify(
            self._get_feature(data)
        )

    def _get_feature(self, datapt):
        tagged_tokens = nltk.tag.pos_tag(
            [token for token in nltk.word_tokenize(datapt)
                if token not in string.punctuation
                and token.lower() not in nltk.corpus.stopwords.words('english')]
        )
        return {self._lemmatizer.lemmatize(token, POS_TAG_MAP[pos[0].upper()]): True
                    for token, pos in tagged_tokens}

    def _get_feature_set(self, data):
        return [self._get_feature(pt) for pt in data]

def main():

    model = NTLKNaiveBayesClassifier()
    pl = Pipeline(model, nltk.corpus.movie_reviews)
    pprint(pl.cross_validate_kfold(10))

if __name__ == '__main__':
    main()
```