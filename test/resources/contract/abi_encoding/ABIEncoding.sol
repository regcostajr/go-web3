// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ABIEncoding {

    uint[] public uint_dynamic_array;

    constructor(uint[] memory _uint_dynamic_array){
        uint_dynamic_array = _uint_dynamic_array;
    }

    event TestString(string _string, string[] _string_dynamic_array, string[2] _string_fixed_array);
    event TestBytes(bytes _bytes, bytes[] _bytes_dynamic_array, bytes[2] _bytes_fixed_array, bytes10[2] _bytes_fixed_both);

    function testString(string memory _string, string[] memory _string_dynamic_array, string[2] memory _string_fixed_array) public {
        emit TestString(_string, _string_dynamic_array, _string_fixed_array);
    }

    function testBytes (bytes memory _bytes, bytes[] memory _bytes_dynamic_array, bytes[2] memory _bytes_fixed_array, bytes10[2] memory _bytes_fixed_both) public {
        emit TestBytes(_bytes, _bytes_dynamic_array, _bytes_fixed_array, _bytes_fixed_both);
    }

}
